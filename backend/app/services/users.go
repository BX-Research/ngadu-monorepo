package services

import (
	"context"
	"errors"
	"log"

	"github.com/BX-Research/ngadu-monorepo/backend/app/types"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

type UserFinder interface {
	FindUser(context.Context, string) (types.UserInfo, error)
}

func FindUser(users UserFinder) fiber.Handler {
	return func(c fiber.Ctx) error {
		user, err := users.FindUser(c.Context(), c.Params("id"))
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(404).JSON(fiber.Map{"error": "user not found"})
		}
		if err != nil {
			log.Printf("user lookup failed: %v", err)
			return c.SendStatus(500)
		}
		return c.JSON(&user)
	}
}
