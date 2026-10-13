package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
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
	return conn
}

func TestSeedRerunPreservesProfilesAndSyncsRoles(t *testing.T) {
	conn := testDB(t)
	ctx := context.Background()
	var output bytes.Buffer
	if err := seed(ctx, conn, &output); err != nil {
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
	if err := seed(ctx, conn, &output); err != nil {
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
	if _, err := conn.Exec(ctx, schemaSQL); err != nil {
		t.Fatal(err)
	}
	// Reserve a fixture ID for an unrelated user, forcing a conflict mid-seed.
	if _, err := conn.Exec(ctx, `INSERT INTO users (id, firebase_uid, display_name)
		VALUES ('00000000-0000-4000-8000-000000000003', 'unrelated-user', 'Existing User')`); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := seed(ctx, conn, &output); err == nil {
		t.Fatal("expected a conflicting ID to fail the seed")
	}
	var count int
	if err := conn.QueryRow(ctx, "SELECT count(*) FROM users").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 || output.Len() != 0 {
		t.Fatalf("failed seed left partial data or success output: users=%d, output=%q", count, output.String())
	}
}
