package routes

import (
	"github.com/BX-Research/ngadu-monorepo/backend/app/services"
	"github.com/gofiber/fiber/v3"
)

// This unauthenticated read API is exclusively for fictional local development.
func LeaderboardRoutes(app fiber.Router, db services.LeaderboardReader, photos services.PhotoSigner) {
	app.Get("/api/leaderboard", services.GetLeaderboard(db, photos))
}
