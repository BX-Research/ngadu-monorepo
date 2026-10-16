package main

import (
	"context"
	"errors"
	"os"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v3"
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
func registerUserRoutes(app *fiber.App, lookup func(context.Context, string) (userInfo, error)) {
	app.Get("/api/users/:id", func(c fiber.Ctx) error {
		c.Set("Cache-Control", "no-store")
		fail := func(code int, message string) error {
			return c.Status(code).JSON(fiber.Map{"error": message}, fiber.MIMEApplicationJSON)
		}
		id := c.Params("id")
		if !userIDPattern.MatchString(id) {
			return fail(fiber.StatusBadRequest, "invalid user ID")
		}
		ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
		defer cancel()
		user, err := lookup(ctx, id)
		if errors.Is(err, pgx.ErrNoRows) {
			return fail(fiber.StatusNotFound, "user not found")
		}
		if err != nil {
			return fail(fiber.StatusServiceUnavailable, "user lookup unavailable")
		}
		return c.JSON(user, fiber.MIMEApplicationJSON)
	})
	app.All("/api/users/:id", methodNotAllowed)
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
