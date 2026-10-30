package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BX-Research/ngadu-monorepo/backend/app/types"
	"github.com/gofiber/fiber/v3"
)

type ladderStub struct {
	err   error
	calls int
}

func (s *ladderStub) Leaderboard(ctx context.Context) (types.Leaderboard, error) {
	s.calls++
	if _, ok := ctx.Deadline(); !ok {
		return types.Leaderboard{}, errors.New("missing deadline")
	}
	return types.Leaderboard{TotalPlayers: 1,
		Players: []types.RankedPlayer{{Rank: 1, UserID: "player", DisplayName: "Test Player", Rating: 1500}}}, s.err
}

type photoStub struct{ err error }

func (s photoStub) PhotoURL(context.Context, string) (string, error) {
	return "http://127.0.0.1:8333/signed", s.err
}

func TestLeaderboardHandler(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{
		{"found", nil, 200},
		{"unavailable", errors.New("secret database details"), 503},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db := &ladderStub{err: tc.err}
			app := fiber.New()
			app.Get("/api/leaderboard", GetLeaderboard(db, photoStub{}))
			res, err := app.Test(httptest.NewRequest("GET", "/api/leaderboard", nil))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != tc.status {
				t.Fatalf("status %d want %d", res.StatusCode, tc.status)
			}
			if res.Header.Get("Cache-Control") != "no-store" {
				t.Fatal("response is cacheable")
			}
			var body map[string]any
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			data, _ := json.Marshal(body)
			if strings.Contains(string(data), "secret") || strings.Contains(string(data), "photo_key") || strings.Contains(string(data), `"community"`) {
				t.Fatal("private details leaked")
			}
			if tc.status == 200 {
				players := body["players"].([]any)
				if players[0].(map[string]any)["photo_url"] == nil {
					t.Fatal("missing signed photo")
				}
			}
		})
	}
}
func TestLeaderboardPhotoFailureDoesNotHideStandings(t *testing.T) {
	app := fiber.New()
	app.Get("/api/leaderboard", GetLeaderboard(&ladderStub{}, photoStub{errors.New("storage unavailable")}))
	res, err := app.Test(httptest.NewRequest("GET", "/api/leaderboard", nil))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var body types.Leaderboard
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 || len(body.Players) != 1 || body.Players[0].PhotoURL != nil {
		t.Fatal("storage failure affected standings")
	}
}
