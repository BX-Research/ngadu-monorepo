package routes

import (
	"github.com/BX-Research/ngadu-monorepo/backend/app/services"
	"github.com/gofiber/fiber/v3"
	"os"
)

// This unauthenticated read API is exclusively for fictional local development.
func LeaderboardRoutes(app fiber.Router, db services.LeaderboardReader, photos services.PhotoSigner) {
	if os.Getenv("ENABLE_DEV_LEADERBOARD") != "true" {
		return
	}
	app.Get("/api/communities/:communityId/leaderboard", services.GetLeaderboard(db, photos))
}
