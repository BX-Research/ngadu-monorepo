package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

// LoadEnv reads the root .env without overriding shell settings.
func LoadEnv() error {
	root, err := Root()
	if err != nil {
		return err
	}
	err = godotenv.Load(filepath.Join(root, ".env"))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func required(key string) (string, error) {
	value := os.Getenv(key)
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return value, nil
}

// URL reads the explicitly configured application database URL.
func URL() (string, error) {
	if err := LoadEnv(); err != nil {
		return "", err
	}
	return required("DATABASE_URL")
}

// LocalURL constructs a loopback URL for Make development commands.
// Ambient DATABASE_URL is deliberately ignored.
func LocalURL() (string, error) {
	if err := LoadEnv(); err != nil {
		return "", err
	}
	user, err := required("DEV_DB_USER")
	if err != nil {
		return "", err
	}
	password, err := required("DEV_DB_PASSWORD")
	if err != nil {
		return "", err
	}
	name, err := required("DEV_DB_NAME")
	if err != nil {
		return "", err
	}
	port := os.Getenv("DEV_DB_PORT")
	if port == "" {
		port = "5432"
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return "", fmt.Errorf("invalid DEV_DB_PORT")
	}
	u := &url.URL{Scheme: "postgresql", Host: net.JoinHostPort("127.0.0.1", port),
		User: url.UserPassword(user, password), Path: "/" + name, RawQuery: "sslmode=disable"}
	return u.String(), nil
}

// Open creates a shared pool for the supplied database URL.
// The caller owns the pool and must close it during shutdown.
func Open(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	if strings.TrimSpace(databaseURL) == "" {
		return nil, fmt.Errorf("database URL is required")
	}
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("connect database: %w", err)
	}
	return pool, nil
}
