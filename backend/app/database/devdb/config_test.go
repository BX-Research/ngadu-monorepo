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
	raw, err := URL()
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
