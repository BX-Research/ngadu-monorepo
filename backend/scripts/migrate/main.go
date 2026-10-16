package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BX-Research/ngadu-monorepo/backend/internal/devdb"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: migrate up|status|create NAME")
	}
	root, err := devdb.Root()
	if err != nil {
		return err
	}
	dir := filepath.Join(root, "backend", "migrations")
	if args[0] == "create" {
		if len(args) != 2 || filepath.Base(args[1]) != args[1] || args[1] == "." || args[1] == ".." {
			return fmt.Errorf("provide a migration name without path separators")
		}
		return goose.Create(nil, dir, args[1], "sql")
	}
	if len(args) != 1 || (args[0] != "up" && args[0] != "status") {
		return fmt.Errorf("supported commands: up, status, create NAME")
	}
	connectionURL, err := devdb.URL()
	if err != nil {
		return err
	}
	db, err := sql.Open("pgx", connectionURL)
	if err != nil {
		return fmt.Errorf("open local development database failed")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("local development database unavailable; run make db-up")
	}
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.RunContext(ctx, args[0], db, dir)
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Migration failed:", err)
		os.Exit(1)
	}
}
