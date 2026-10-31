package routes

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/BX-Research/ngadu-monorepo/backend/app/types"
	"github.com/gofiber/fiber/v3"
)

type leaderboardReader struct{}

func (leaderboardReader) Leaderboard(_ context.Context) (types.Leaderboard, error) {
	return types.Leaderboard{Players: []types.RankedPlayer{}}, nil
}

func TestLeaderboardAvailableByDefault(t *testing.T) {
	app := fiber.New()
	LeaderboardRoutes(app, leaderboardReader{}, nil)
	res, err := app.Test(httptest.NewRequest("GET", "/api/leaderboard", nil))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("status %d, want 200", res.StatusCode)
	}
}

func TestScopedLeaderboardRouteRemoved(t *testing.T) {
	app := fiber.New()
	LeaderboardRoutes(app, leaderboardReader{}, nil)
	res, err := app.Test(httptest.NewRequest("GET", "/api/communities/00000000-0000-4000-8000-000000000101/leaderboard", nil))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 404 {
		t.Fatalf("old route status %d, want 404", res.StatusCode)
	}
}
