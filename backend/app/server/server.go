// Package server builds the HTTP application from injected dependencies.
package server

import (
	"context"
	"time"

	"github.com/BX-Research/ngadu-monorepo/backend/app/routes"
	"github.com/BX-Research/ngadu-monorepo/backend/app/services"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	Users       services.UserFinder
	Leaderboard services.LeaderboardReader
	Photos      services.PhotoSigner
	Ready       func(context.Context) error
}

// New never loads configuration or constructs database/storage clients.
func New(deps Dependencies) *fiber.App {
	app := fiber.New(fiber.Config{CaseSensitive: true, ReadTimeout: 5 * time.Second, IdleTimeout: 60 * time.Second})
	app.Get("/healthz", func(c fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) })
	app.Get("/readyz", func(c fiber.Ctx) error {
		c.Set("Cache-Control", "no-store")
		ctx, cancel := context.WithTimeout(c.Context(), time.Second)
		defer cancel()
		if deps.Ready == nil || deps.Ready(ctx) != nil {
			return c.Status(503).JSON(fiber.Map{"status": "unavailable"})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})
	routes.LeaderboardRoutes(app, deps.Leaderboard, deps.Photos)
	routes.UserRoutes(app, deps.Users)
	return app
}
