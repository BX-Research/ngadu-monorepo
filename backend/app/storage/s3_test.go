package storage

import (
	"context"
	"net/url"
	"strings"
	"testing"
)

func TestSignedPhotoURL(t *testing.T) {
	store, err := New(Config{Endpoint: "http://127.0.0.1:8333", Region: "us-east-1", Bucket: "ngadu-dev", AccessKey: "local-key", SecretKey: "local-secret"})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := store.PhotoURL(context.Background(), "00000000-0000-4000-8000-000000000201")
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if u.Host != "127.0.0.1:8333" || u.Path != "/ngadu-dev/avatars/00000000-0000-4000-8000-000000000201.jpg" || u.Query().Get("X-Amz-Expires") != "900" || u.Query().Get("X-Amz-Signature") == "" {
		t.Fatalf("bad signed URL: %s", raw)
	}
	if strings.Contains(raw, "local-secret") {
		t.Fatal("secret exposed")
	}
}
func TestStorageConfigValidation(t *testing.T) {
	for _, endpoint := range []string{"", ":bad", "file:///tmp/bucket", "http://user:password@127.0.0.1:8333", "http://127.0.0.1:8333?x=1"} {
		if _, err := New(Config{Endpoint: endpoint, Region: "us-east-1", Bucket: "ngadu-dev", AccessKey: "key", SecretKey: "secret"}); err == nil {
			t.Fatalf("accepted %q", endpoint)
		}
	}
}

func TestConfiguredEndpointIsNotRestrictedToDevelopment(t *testing.T) {
	cfg := Config{EndpointHost: "photos.example", Endpoint: "https://photos.example", Region: "region", Bucket: "profiles", AccessKey: "key", SecretKey: "secret"}
	store, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := store.PhotoURL(context.Background(), "00000000-0000-4000-8000-000000000201")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(raw, "https://photos.example/profiles/avatars/00000000-0000-4000-8000-000000000201.jpg?") {
		t.Fatal("configured endpoint was ignored")
	}
	if err := cfg.ValidateLocal(); err == nil {
		t.Fatal("development seed accepted remote endpoint")
	}
	cfg.EndpointHost = "different.example"
	if _, err := New(cfg); err == nil {
		t.Fatal("mismatched endpoint host accepted")
	}
	cfg.EndpointHost = "localhost"
	cfg.Endpoint = "http://localhost:9000"
	if err := cfg.ValidateLocal(); err != nil {
		t.Fatal(err)
	}
}
