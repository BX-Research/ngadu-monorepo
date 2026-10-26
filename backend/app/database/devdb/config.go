// Package devdb configures commands that must target local development PostgreSQL.
package devdb

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Root supports invoking commands from the repository root or backend directory.
func Root() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "backend", "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("run from the Ngadu repository")
		}
		dir = parent
	}
}

func settings() (func(string, string) string, error) {
	root, err := Root()
	if err != nil {
		return nil, err
	}
	values, err := godotenv.Read(filepath.Join(root, ".env"))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("read development .env: %w", err)
	}
	setting := func(key, fallback string) string {
		if value, exists := os.LookupEnv(key); exists {
			if value == "" {
				return fallback
			}
			return value
		}
		if value := values[key]; value != "" {
			return value
		}
		return fallback
	}
	return setting, nil
}

// URL constructs a local database URL, ignoring DATABASE_URL.
func URL() (string, error) {
	setting, err := settings()
	if err != nil {
		return "", err
	}
	u := &url.URL{Scheme: "postgresql", Host: "127.0.0.1:5432",
		User: url.UserPassword(setting("DEV_DB_USER", "ngadu"), setting("DEV_DB_PASSWORD", "change-me-local-only")),
		Path: "/" + setting("DEV_DB_NAME", "ngadu_dev"), RawQuery: "sslmode=disable"}
	return u.String(), nil
}
