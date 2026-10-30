package dal

import (
	"context"
	"errors"
	"github.com/BX-Research/ngadu-monorepo/backend/app/database"
	"github.com/BX-Research/ngadu-monorepo/backend/internal/testdb"
	"github.com/jackc/pgx/v5"
	"os"
	"path/filepath"
	"testing"
)

func TestLeaderboardQuery(t *testing.T) {
	db := testdb.Open(t)
	ctx := context.Background()
	root, err := database.Root()
	if err != nil {
		t.Fatal(err)
	}
	seed, err := os.ReadFile(filepath.Join(root, "backend/seeds/dev.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(ctx, string(seed)); err != nil {
		t.Fatal(err)
	}
	ladder := NewLeaderboard(db)
	got, err := ladder.Leaderboard(ctx, "00000000-0000-4000-8000-000000000101")
	if err != nil {
		t.Fatal(err)
	}
	if got.TotalPlayers != 23 || len(got.Players) != 10 {
		t.Fatalf("unexpected count: %+v", got)
	}
	if got.Players[1].Rating != 1690 || got.Players[2].Rating != 1690 || got.Players[1].UserID >= got.Players[2].UserID {
		t.Fatal("unstable ties")
	}
	for i, p := range got.Players {
		if p.Rank != i+1 {
			t.Fatal("incorrect rank")
		}
		if i > 0 && p.Rating > got.Players[i-1].Rating {
			t.Fatal("out of order")
		}
	}
	other, err := ladder.Leaderboard(ctx, "00000000-0000-4000-8000-000000000102")
	if err != nil {
		t.Fatal(err)
	}
	if other.TotalPlayers != 2 || len(other.Players) != 2 || other.Players[0].Rating != 1280 {
		t.Fatal("community isolation failed")
	}
	if _, err = db.Exec(ctx, `UPDATE player_ratings SET rating=2500 WHERE community_id='00000000-0000-4000-8000-000000000101' AND user_id='00000000-0000-4000-8000-000000000220'`); err != nil {
		t.Fatal(err)
	}
	refreshed, err := ladder.Leaderboard(ctx, got.Community.ID)
	if err != nil || refreshed.Players[0].Rating != 2500 {
		t.Fatal("rating edit not reflected")
	}
	if _, err = db.Exec(ctx, `DELETE FROM player_ratings WHERE community_id='00000000-0000-4000-8000-000000000102'`); err != nil {
		t.Fatal(err)
	}
	empty, err := ladder.Leaderboard(ctx, other.Community.ID)
	if err != nil || len(empty.Players) != 0 || empty.TotalPlayers != 0 || empty.Players == nil {
		t.Fatalf("empty: %+v %v", empty, err)
	}
	_, err = ladder.Leaderboard(ctx, "00000000-0000-4000-8000-000000000999")
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("unknown: %v", err)
	}
}
