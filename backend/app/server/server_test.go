package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/BX-Research/ngadu-monorepo/backend/app/types"
	"net/http/httptest"
	"testing"
)

type ladder struct{}

func (ladder) Leaderboard(context.Context) (types.Leaderboard, error) {
	return types.Leaderboard{TotalPlayers: 1, Players: []types.RankedPlayer{{UserID: "player", Rank: 1}}}, nil
}

type photos struct{ calls int }

func (p *photos) PhotoURL(_ context.Context, key string) (string, error) {
	p.calls++
	return "https://photos.example/" + key, nil
}

func TestNewUsesInjectedDependencies(t *testing.T) {
	signer := &photos{}
	app := New(Dependencies{Leaderboard: ladder{}, Photos: signer, Ready: func(context.Context) error { return errors.New("offline") }})
	response, err := app.Test(httptest.NewRequest("GET", "/api/leaderboard", nil))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var result types.Leaderboard
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 200 || signer.calls != 1 || result.Players[0].PhotoURL == nil || *result.Players[0].PhotoURL != "https://photos.example/player" {
		t.Fatal("photo signer was not injected")
	}
	for _, tc := range []struct {
		path   string
		status int
	}{{"/healthz", 200}, {"/readyz", 503}} {
		res, err := app.Test(httptest.NewRequest("GET", tc.path, nil))
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != tc.status {
			t.Fatalf("%s: %d", tc.path, res.StatusCode)
		}
	}
}

func TestNewWithoutPhotosKeepsLeaderboardAvailable(t *testing.T) {
	app := New(Dependencies{Leaderboard: ladder{}, Ready: func(context.Context) error { return nil }})
	response, err := app.Test(httptest.NewRequest("GET", "/api/leaderboard", nil))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var result types.Leaderboard
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 200 || len(result.Players) != 1 || result.Players[0].PhotoURL != nil {
		t.Fatal("optional photo dependency affected standings")
	}
	ready, err := app.Test(httptest.NewRequest("GET", "/readyz", nil))
	if err != nil {
		t.Fatal(err)
	}
	defer ready.Body.Close()
	if ready.StatusCode != 200 {
		t.Fatal("injected readiness check was not used")
	}
}
