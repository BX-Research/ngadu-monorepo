package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

func testDB(t *testing.T) *pgx.Conn {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set TEST_DATABASE_URL to run PostgreSQL integration tests")
	}
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	schema := pgx.Identifier{fmt.Sprintf("seed_test_%d", time.Now().UnixNano())}.Sanitize()
	if _, err := conn.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		conn.Close(ctx)
		t.Fatal(err)
	}
	t.Cleanup(func() {
		conn.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE")
		conn.Close(ctx)
	})
	if _, err := conn.Exec(ctx, "SET search_path TO "+schema); err != nil {
		t.Fatal(err)
	}
	migration, err := os.ReadFile("../migrations/20261013000000_initial_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	up, _, _ := strings.Cut(string(migration), "-- +goose Down")
	if _, err := conn.Exec(ctx, up); err != nil {
		t.Fatal(err)
	}
	return conn
}

func TestSeedRerunPreservesProfilesAndSyncsRoles(t *testing.T) {
	conn := testDB(t)
	ctx := context.Background()

	if err := seedSQL(ctx, conn); err != nil {
		t.Fatal(err)
	}
	// Existing Firebase users can have IDs assigned outside the seed script.
	_, err := conn.Exec(ctx, `
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
	if users != 4 || communities != 2 || memberships != 5 || name != "Edited Demo" || role != "organizer" {
		t.Fatalf("unexpected seeded state: %d users, %d communities, %d memberships, name=%q role=%q", users, communities, memberships, name, role)
	}
}

func TestSeedRollsBackOnConflict(t *testing.T) {
	conn := testDB(t)
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
