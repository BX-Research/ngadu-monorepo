package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// Run after make storage-seed against an isolated local S3 service.
func TestPrivateFixturePhotos(t *testing.T) {
	endpoint := os.Getenv("TEST_S3_ENDPOINT")
	if endpoint == "" {
		t.Skip("set TEST_S3_ENDPOINT after seeding disposable local storage")
	}
	cfg := DevConfig()
	cfg.Endpoint = endpoint
	if err := cfg.ValidateLocal(); err != nil {
		t.Fatal(err)
	}
	store, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	for _, tc := range []struct {
		key    string
		signed bool
		status int
	}{
		{"profiles/sage.svg", true, 200},
		{"profiles/sage.svg", false, 403},
		{"profiles/missing.svg", true, 404},
	} {
		t.Run(fmt.Sprintf("%s/signed=%t", tc.key, tc.signed), func(t *testing.T) {
			raw := cfg.Endpoint + "/" + cfg.Bucket + "/" + tc.key
			if tc.signed {
				raw, err = store.PhotoURL(context.Background(), tc.key)
				if err != nil {
					t.Fatal(err)
				}
			}
			res, err := client.Get(raw)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != tc.status {
				t.Fatalf("status %d want %d", res.StatusCode, tc.status)
			}
			if tc.status == 200 {
				body, _ := io.ReadAll(res.Body)
				if !strings.Contains(string(body), "<svg") {
					t.Fatal("fixture is not an image")
				}
			}
		})
	}
}
