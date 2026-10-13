package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5"
)

type userInfo struct {
	ID          string    `json:"id"`
	FirebaseUID string    `json:"firebase_uid"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
}

var userIDPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// registerUserRoutes exposes local development profile inspection. Authentication
// and community authorization must be added before this becomes a product API.
func registerUserRoutes(mux *http.ServeMux, lookup func(context.Context, string) (userInfo, error)) {
	mux.HandleFunc("GET /api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		fail := func(code int, message string) {
			w.WriteHeader(code)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
		}
		id := r.PathValue("id")
		if !userIDPattern.MatchString(id) {
			fail(http.StatusBadRequest, "invalid user ID")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		user, err := lookup(ctx, id)
		if errors.Is(err, pgx.ErrNoRows) {
			fail(http.StatusNotFound, "user not found")
			return
		}
		if err != nil {
			fail(http.StatusServiceUnavailable, "user lookup unavailable")
			return
		}
		_ = json.NewEncoder(w).Encode(user)
	})
}

func findUser(ctx context.Context, id string) (userInfo, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return userInfo{}, err
	}
	defer conn.Close(context.Background())
	var user userInfo
	err = conn.QueryRow(ctx, `SELECT id::text, firebase_uid, display_name, created_at FROM users WHERE id = $1`, id).
		Scan(&user.ID, &user.FirebaseUID, &user.DisplayName, &user.CreatedAt)
	return user, err
}
