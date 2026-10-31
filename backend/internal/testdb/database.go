// Package testdb provides isolated schemas in an explicitly selected disposable database.
package testdb

import (
	"context"
	"fmt"
	"github.com/BX-Research/ngadu-monorepo/backend/app/database"
	"github.com/jackc/pgx/v5"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func Open(t *testing.T) *pgx.Conn {
	return OpenBefore(t, "")
}

// OpenBefore stops before a named migration to exercise upgrades from old schemas.
func OpenBefore(t *testing.T, stopBefore string) *pgx.Conn {
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
	root, err := database.Root()
	if err != nil {
		t.Fatal(err)
	}
	paths, err := filepath.Glob(filepath.Join(root, "backend", "migrations", "*.sql"))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		if filepath.Base(path) == stopBefore {
			break
		}
		migration, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		up, _, _ := strings.Cut(string(migration), "-- +goose Down")
		if _, err := conn.Exec(ctx, up); err != nil {
			t.Fatal(err)
		}
	}
	return conn
}
