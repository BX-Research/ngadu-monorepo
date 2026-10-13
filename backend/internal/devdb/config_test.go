package devdb

import (
	"net/url"
	"testing"
)

func TestLocalURLIgnoresAmbientDatabase(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://remote.invalid/production")
	t.Setenv("DEV_DB_USER", "user@example")
	t.Setenv("DEV_DB_PASSWORD", "p@ss:/?#")
	t.Setenv("DEV_DB_NAME", "dev database")
	t.Setenv("DEV_DB_PORT", "55432")
	raw, err := URL()
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	password, _ := u.User.Password()
	if u.Host != "127.0.0.1:55432" || u.Path != "/dev database" || u.User.Username() != "user@example" || password != "p@ss:/?#" {
		t.Fatalf("unexpected local configuration")
	}
}

func TestInvalidPort(t *testing.T) {
	t.Setenv("DEV_DB_PORT", "0")
	if _, err := URL(); err == nil {
		t.Fatal("expected invalid port to fail")
	}
}

func TestAPIAddressStaysLocal(t *testing.T) {
	t.Setenv("API_ADDR", "remote.invalid:8080")
	t.Setenv("API_PORT", "18080")
	addr, err := APIAddress()
	if err != nil {
		t.Fatal(err)
	}
	if addr != "127.0.0.1:18080" {
		t.Fatalf("unexpected address %q", addr)
	}
	t.Setenv("API_PORT", "invalid")
	if _, err := APIAddress(); err == nil {
		t.Fatal("expected invalid API port to fail")
	}
}
