package database

import (
	"net/url"
	"testing"
)

func TestLocalURLIgnoresAmbientDatabase(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://remote.invalid/production")
	t.Setenv("DEV_DB_USER", "user@example")
	t.Setenv("DEV_DB_PASSWORD", "p@ss:/?#")
	t.Setenv("DEV_DB_NAME", "dev database")
	raw, err := LocalURL()
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	password, _ := u.User.Password()
	if u.Host != "127.0.0.1:5432" || u.Path != "/dev database" || u.User.Username() != "user@example" || password != "p@ss:/?#" {
		t.Fatalf("unexpected local configuration")
	}
}

func TestLocalURLRequiresSettings(t *testing.T) {
	for _, key := range []string{"DEV_DB_USER", "DEV_DB_PASSWORD", "DEV_DB_NAME"} {
		t.Run(key, func(t *testing.T) {
			t.Setenv("DEV_DB_USER", "user")
			t.Setenv("DEV_DB_PASSWORD", "password")
			t.Setenv("DEV_DB_NAME", "db")
			t.Setenv(key, "")
			if _, err := LocalURL(); err == nil {
				t.Fatal("expected missing setting error")
			}
		})
	}
}
func TestURLRequiresExplicitValue(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	if _, err := URL(); err == nil {
		t.Fatal("expected missing URL error")
	}
	t.Setenv("DATABASE_URL", "postgres://example.invalid/test")
	got, err := URL()
	if err != nil || got != "postgres://example.invalid/test" {
		t.Fatalf("unexpected URL result: %v", err)
	}
}
