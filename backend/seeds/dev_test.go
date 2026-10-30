package main

import (
	"context"

	"os"

	"testing"

	"github.com/BX-Research/ngadu-monorepo/backend/internal/testdb"
	"github.com/jackc/pgx/v5"
)

func TestSeedRerunPreservesProfilesAndSyncsRoles(t *testing.T) {
	conn := testdb.Open(t)
	ctx := context.Background()

	if err := seedSQL(ctx, conn); err != nil {
		t.Fatal(err)
	}
	// Existing Firebase users can have IDs assigned outside the seed script.
	_, err := conn.Exec(ctx, `
		DELETE FROM player_ratings WHERE user_id = '00000000-0000-4000-8000-000000000002';
		UPDATE users SET display_name = 'Edited Demo' WHERE firebase_uid = 'seed-user-avery';
		UPDATE memberships SET role = 'member' WHERE user_id = '00000000-0000-4000-8000-000000000001';
		DELETE FROM memberships WHERE user_id = '00000000-0000-4000-8000-000000000002';
		UPDATE users SET id = '00000000-0000-4000-8000-000000000099' WHERE firebase_uid = 'seed-user-blake';
	`)
	if err != nil {
		t.Fatal(err)
	}
	if err := seedSQL(ctx, conn); err != nil {
		t.Fatal(err)
	}
	var users, communities, memberships int
	var name, role string
	err = conn.QueryRow(ctx, `SELECT
		(SELECT count(*) FROM users),
		(SELECT count(*) FROM communities),
		(SELECT count(*) FROM memberships),
		(SELECT display_name FROM users WHERE firebase_uid = 'seed-user-avery'),
		(SELECT role FROM memberships WHERE user_id = '00000000-0000-4000-8000-000000000001'
		 AND community_id = '00000000-0000-4000-8000-000000000101')
	`).Scan(&users, &communities, &memberships, &name, &role)
	if err != nil {
		t.Fatal(err)
	}
	if users != 24 || communities != 2 || memberships != 25 || name != "Edited Demo" || role != "organizer" {
		t.Fatalf("unexpected seeded state: %d users, %d communities, %d memberships, name=%q role=%q", users, communities, memberships, name, role)
	}
}

func TestSeedRollsBackOnConflict(t *testing.T) {
	conn := testdb.Open(t)
	ctx := context.Background()
	// Reserve a fixture ID for an unrelated user, forcing a conflict mid-seed.
	if _, err := conn.Exec(ctx, `INSERT INTO users (id, firebase_uid, display_name)
		VALUES ('00000000-0000-4000-8000-000000000003', 'unrelated-user', 'Existing User')`); err != nil {
		t.Fatal(err)
	}

	if err := seedSQL(ctx, conn); err == nil {
		t.Fatal("expected a conflicting ID to fail the seed")
	}
	var count int
	if err := conn.QueryRow(ctx, "SELECT count(*) FROM users").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("failed seed left partial data: users=%d", count)
	}
}

func seedSQL(ctx context.Context, conn *pgx.Conn) error {
	data, err := os.ReadFile("dev.sql")
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, string(data))
	if err != nil {
		conn.Exec(ctx, "ROLLBACK")
	}
	return err
}

func TestLeaderboardFixturesPreserveEdits(t *testing.T) {
	conn := testdb.Open(t)
	ctx := context.Background()
	if err := seedSQL(ctx, conn); err != nil {
		t.Fatal(err)
	}
	var count, missing, distinctRatings int
	if err := conn.QueryRow(ctx, `SELECT count(*), count(*) FILTER (WHERE u.profile_photo_key IS NULL), count(DISTINCT r.rating)
 FROM player_ratings r JOIN users u ON u.id=r.user_id
 WHERE r.community_id='00000000-0000-4000-8000-000000000101'`).Scan(&count, &missing, &distinctRatings); err != nil {
		t.Fatal(err)
	}
	if count != 23 || missing == 0 || distinctRatings >= count {
		t.Fatalf("fixtures: count=%d missing=%d ratings=%d", count, missing, distinctRatings)
	}
	if _, err := conn.Exec(ctx, `UPDATE player_ratings SET rating=2222, wins=99 WHERE user_id='00000000-0000-4000-8000-000000000201';
 UPDATE users SET profile_photo_key='custom/photo.png' WHERE id='00000000-0000-4000-8000-000000000201'`); err != nil {
		t.Fatal(err)
	}
	if err := seedSQL(ctx, conn); err != nil {
		t.Fatal(err)
	}
	var rating, wins int
	var key string
	if err := conn.QueryRow(ctx, `SELECT rating,wins,profile_photo_key FROM player_ratings JOIN users ON users.id=user_id
 WHERE user_id='00000000-0000-4000-8000-000000000201'`).Scan(&rating, &wins, &key); err != nil {
		t.Fatal(err)
	}
	if rating != 2222 || wins != 99 || key != "custom/photo.png" {
		t.Fatal("seed overwrote edits")
	}
	if _, err := conn.Exec(ctx, `INSERT INTO player_ratings(community_id,user_id,rating,wins,losses)
 VALUES('00000000-0000-4000-8000-000000000101','00000000-0000-4000-8000-000000000004',1200,0,0)`); err == nil {
		t.Fatal("nonmember received rating")
	}
	if _, err := conn.Exec(ctx, `UPDATE player_ratings SET wins=-1`); err == nil {
		t.Fatal("negative wins accepted")
	}
}
