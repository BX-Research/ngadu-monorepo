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
	raw, err := store.PhotoURL(context.Background(), "profiles/test person.svg")
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if u.Host != "127.0.0.1:8333" || u.Path != "/ngadu-dev/profiles/test person.svg" || u.Query().Get("X-Amz-Expires") != "900" || u.Query().Get("X-Amz-Signature") == "" {
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
