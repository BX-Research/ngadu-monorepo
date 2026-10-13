package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

func TestUserRoute(t *testing.T) {
	const id = "00000000-0000-4000-8000-000000000001"
	created := time.Date(2026, 10, 13, 0, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name, method, path string
		err                error
		status             int
		called             bool
	}{
		{"found", "GET", "/api/users/" + id, nil, 200, true},
		{"invalid ID", "GET", "/api/users/not-a-uuid", nil, 400, false},
		{"missing user", "GET", "/api/users/" + id, pgx.ErrNoRows, 404, true},
		{"database unavailable", "GET", "/api/users/" + id, errors.New("private database credentials"), 503, true},
		{"unsupported method", "POST", "/api/users/" + id, nil, 405, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			called := false
			registerUserRoutes(mux, func(ctx context.Context, got string) (userInfo, error) {
				called = true
				if got != id {
					t.Fatalf("lookup ID = %q", got)
				}
				if _, ok := ctx.Deadline(); !ok {
					t.Fatal("lookup requires deadline")
				}
				return userInfo{ID: id, FirebaseUID: "seed-user-avery", DisplayName: "Avery Demo", CreatedAt: created}, tc.err
			})
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
			if w.Code != tc.status || called != tc.called {
				t.Fatalf("status=%d called=%v body=%s", w.Code, called, w.Body.String())
			}
			if strings.Contains(w.Body.String(), "private") {
				t.Fatal("database details exposed")
			}
			if tc.status == 200 {
				var got userInfo
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatal(err)
				}
				if got.ID != id || got.FirebaseUID != "seed-user-avery" || got.DisplayName != "Avery Demo" || !got.CreatedAt.Equal(created) {
					t.Fatalf("unexpected user: %+v", got)
				}
			}
			if tc.status != 405 && w.Header().Get("Content-Type") != "application/json" {
				t.Fatal("expected JSON response")
			}
		})
	}
}

func TestUserRouteRegistered(t *testing.T) {
	w := httptest.NewRecorder()
	newHandler(func(context.Context) error { return nil }).ServeHTTP(w, httptest.NewRequest("GET", "/api/users/invalid", nil))
	if w.Code != 400 {
		t.Fatalf("route missing: status %d", w.Code)
	}
}

func TestFindUserPostgres(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("set TEST_DATABASE_URL to run PostgreSQL integration tests")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, raw)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(context.Background())
	schema := fmt.Sprintf("user_lookup_test_%d", time.Now().UnixNano())
	quoted := pgx.Identifier{schema}.Sanitize()
	if _, err := conn.Exec(ctx, "CREATE SCHEMA "+quoted); err != nil {
		t.Fatal(err)
	}
	defer conn.Exec(context.Background(), "DROP SCHEMA "+quoted+" CASCADE")
	if _, err := conn.Exec(ctx, "SET search_path TO "+quoted); err != nil {
		t.Fatal(err)
	}
	migration, err := os.ReadFile("../../migrations/20261013000000_initial_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	up, _, _ := strings.Cut(string(migration), "-- +goose Down")
	if _, err := conn.Exec(ctx, up); err != nil {
		t.Fatal(err)
	}
	const id = "00000000-0000-4000-8000-000000000001"
	if _, err := conn.Exec(ctx, "INSERT INTO users (id, firebase_uid, display_name) VALUES ($1, $2, $3)", id, "test-firebase-uid", "Test User"); err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	q.Set("search_path", schema)
	u.RawQuery = q.Encode()
	t.Setenv("DATABASE_URL", u.String())
	user, err := findUser(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if user.ID != id || user.FirebaseUID != "test-firebase-uid" || user.DisplayName != "Test User" || user.CreatedAt.IsZero() {
		t.Fatalf("unexpected user: %+v", user)
	}
	if _, err := findUser(ctx, "00000000-0000-4000-8000-000000000002"); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("missing user error: %v", err)
	}
}
