package services

import (
	"context"
	"github.com/BX-Research/ngadu-monorepo/backend/app/types"
	"github.com/gofiber/fiber/v3"
	"strings"
	"time"
)

type LeaderboardReader interface {
	Leaderboard(context.Context) (types.Leaderboard, error)
}
type PhotoSigner interface {
	PhotoURL(context.Context, string) (string, error)
}

func GetLeaderboard(db LeaderboardReader, photos PhotoSigner) fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set("Cache-Control", "no-store")
		fail := func(status int, message string) error { return c.Status(status).JSON(fiber.Map{"error": message}) }
		ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
		defer cancel()
		result, err := db.Leaderboard(ctx)
		if err != nil {
			return fail(503, "leaderboard unavailable")
		}
		if result.Players == nil {
			result.Players = []types.RankedPlayer{}
		}
		for i := range result.Players {
			p := &result.Players[i]
			p.PhotoURL = nil
			if photos != nil && strings.TrimSpace(p.PhotoKey) != "" {
				if url, err := photos.PhotoURL(ctx, p.PhotoKey); err == nil && url != "" {
					p.PhotoURL = &url
				}
			}
		}
		return c.JSON(result)
	}
}
