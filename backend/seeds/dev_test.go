package seeds

import (
	"context"

	"os"

	"testing"

	"github.com/BX-Research/ngadu-monorepo/backend/internal/testdb"
	"github.com/jackc/pgx/v5"
)

func TestSeedRerunPreservesProfilesAndStoredIDs(t *testing.T) {
	conn := testdb.Open(t)
	ctx := context.Background()

	if err := seedSQL(ctx, conn); err != nil {
		t.Fatal(err)
	}
	// Existing Firebase users can have IDs assigned outside the seed script.
	_, err := conn.Exec(ctx, `
		DELETE FROM player_ratings WHERE user_id = '00000000-0000-4000-8000-000000000002';
		UPDATE users SET display_name = 'Edited Demo' WHERE firebase_uid = 'seed-user-avery';
		UPDATE users SET id = '00000000-0000-4000-8000-000000000099' WHERE firebase_uid = 'seed-user-blake';
	`)
	if err != nil {
		t.Fatal(err)
	}
	if err := seedSQL(ctx, conn); err != nil {
		t.Fatal(err)
	}
	var users, ratings int
	var name, storedID string
	err = conn.QueryRow(ctx, `SELECT
 (SELECT count(*) FROM users), (SELECT count(*) FROM player_ratings),
 (SELECT display_name FROM users WHERE firebase_uid='seed-user-avery'),
 (SELECT user_id::text FROM player_ratings JOIN users ON users.id=user_id WHERE firebase_uid='seed-user-blake')
 `).Scan(&users, &ratings, &name, &storedID)
	if err != nil {
		t.Fatal(err)
	}
	if users != 24 || ratings != 24 || name != "Edited Demo" || storedID != "00000000-0000-4000-8000-000000000099" {
		t.Fatalf("unexpected seeded state: users=%d ratings=%d name=%q id=%q", users, ratings, name, storedID)
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
	var count, distinctRatings int
	if err := conn.QueryRow(ctx, `SELECT count(*), count(DISTINCT r.rating)
 FROM player_ratings r JOIN users u ON u.id=r.user_id`).Scan(&count, &distinctRatings); err != nil {
		t.Fatal(err)
	}
	if count != 24 || distinctRatings >= count {
		t.Fatalf("fixtures: count=%d ratings=%d", count, distinctRatings)
	}
	if _, err := conn.Exec(ctx, `UPDATE player_ratings SET rating=2222, wins=99 WHERE user_id='00000000-0000-4000-8000-000000000201';
 UPDATE users SET display_name='Edited Avatar Demo' WHERE id='00000000-0000-4000-8000-000000000201'`); err != nil {
		t.Fatal(err)
	}
	if err := seedSQL(ctx, conn); err != nil {
		t.Fatal(err)
	}
	var rating, wins int
	var name string
	if err := conn.QueryRow(ctx, `SELECT rating,wins,display_name FROM player_ratings JOIN users ON users.id=user_id
 WHERE user_id='00000000-0000-4000-8000-000000000201'`).Scan(&rating, &wins, &name); err != nil {
		t.Fatal(err)
	}
	if rating != 2222 || wins != 99 || name != "Edited Avatar Demo" {
		t.Fatal("seed overwrote edits")
	}
	if _, err := conn.Exec(ctx, `INSERT INTO player_ratings(user_id,rating,wins,losses) VALUES('00000000-0000-4000-8000-000000000999',1200,0,0)`); err == nil {
		t.Fatal("rating without user accepted")
	}
	if _, err := conn.Exec(ctx, `INSERT INTO player_ratings(user_id,rating,wins,losses) VALUES('00000000-0000-4000-8000-000000000004',1200,0,0)`); err == nil {
		t.Fatal("duplicate player rating accepted")
	}
	if _, err := conn.Exec(ctx, `UPDATE player_ratings SET wins=-1`); err == nil {
		t.Fatal("negative wins accepted")
	}
}
