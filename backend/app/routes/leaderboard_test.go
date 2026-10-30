package routes

import (
	"github.com/gofiber/fiber/v3"
	"net/http/httptest"
	"testing"
)

func TestLeaderboardDisabledByDefault(t *testing.T) {
	t.Setenv("ENABLE_DEV_LEADERBOARD", "")
	app := fiber.New()
	LeaderboardRoutes(app, nil, nil)
	res, err := app.Test(httptest.NewRequest("GET", "/api/communities/00000000-0000-4000-8000-000000000101/leaderboard", nil))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 404 {
		t.Fatalf("disabled status %d", res.StatusCode)
	}
}
