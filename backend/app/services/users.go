package services

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v3"

	"github.com/BX-Research/ngadu-monorepo/backend/app/database/devdb"
	"github.com/BX-Research/ngadu-monorepo/backend/app/types"
	"github.com/jackc/pgx/v5"
)

func FindUser(c fiber.Ctx) error {
	userID := c.Params("id")
	DB_URL, err := devdb.URL()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := pgx.Connect(c, DB_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())
	var user types.UserInfo
	err = conn.QueryRow(c, `SELECT id::text, display_name FROM users WHERE id = $1`, userID).
		Scan(&user.ID, &user.DisplayName)

	if err == pgx.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	if err != nil {
		log.Printf("user lookup failed: %v", err)
		return c.SendStatus(500)
	}

	return c.JSON(&user)
}
