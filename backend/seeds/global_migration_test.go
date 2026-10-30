package main

import (
	"context"
	"github.com/BX-Research/ngadu-monorepo/backend/app/database"
	"github.com/BX-Research/ngadu-monorepo/backend/internal/testdb"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGlobalLadderMigrationResetsRatingsAndPreservesProfiles(t *testing.T) {
	const filename = "20261030000001_global_ladder.sql"
	db := testdb.OpenBefore(t, filename)
	ctx := context.Background()
	_, err := db.Exec(ctx, `INSERT INTO users(id,firebase_uid,display_name,profile_photo_key) VALUES
 ('00000000-0000-4000-8000-000000000001','seed-user-avery','Edited Player','custom/photo.svg'),
 ('00000000-0000-4000-8000-000000000900','migration-unrated','Unrated Player',NULL);
 INSERT INTO communities(id,name) VALUES
 ('00000000-0000-4000-8000-000000000101','First'),('00000000-0000-4000-8000-000000000102','Second');
 INSERT INTO memberships(user_id,community_id,role) VALUES
 ('00000000-0000-4000-8000-000000000001','00000000-0000-4000-8000-000000000101','member'),
 ('00000000-0000-4000-8000-000000000001','00000000-0000-4000-8000-000000000102','organizer');
 INSERT INTO player_ratings(community_id,user_id,rating,wins,losses) VALUES
 ('00000000-0000-4000-8000-000000000101','00000000-0000-4000-8000-000000000001',1800,31,9),
 ('00000000-0000-4000-8000-000000000102','00000000-0000-4000-8000-000000000001',1500,8,3);`)
	if err != nil {
		t.Fatal(err)
	}
	root, err := database.Root()
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "backend/migrations", filename))
	if err != nil {
		t.Fatal(err)
	}
	up, _, _ := strings.Cut(string(data), "-- +goose Down")
	if _, err = db.Exec(ctx, up); err != nil {
		t.Fatal(err)
	}
	var rating, wins, losses, count, obsolete int
	var name, key string
	if err = db.QueryRow(ctx, `SELECT count(*) FROM player_ratings`).Scan(&count); err != nil || count != 0 {
		t.Fatal("old ratings were not reset")
	}
	if err = seedSQL(ctx, db); err != nil {
		t.Fatal(err)
	}
	err = db.QueryRow(ctx, `SELECT rating,wins,losses,display_name,profile_photo_key FROM player_ratings JOIN users ON users.id=user_id WHERE firebase_uid='seed-user-avery'`).Scan(&rating, &wins, &losses, &name, &key)
	if err != nil {
		t.Fatal(err)
	}
	if rating != 1400 || wins != 15 || losses != 8 || name != "Edited Player" || key != "custom/photo.svg" {
		t.Fatal("seed baseline or existing profile changed")
	}
	if err = db.QueryRow(ctx, `SELECT count(*) FROM player_ratings`).Scan(&count); err != nil || count != 24 {
		t.Fatal("global seed ratings missing")
	}
	if err = db.QueryRow(ctx, `SELECT count(*) FROM information_schema.tables WHERE table_schema=current_schema() AND table_name IN ('communities','memberships')`).Scan(&obsolete); err != nil || obsolete != 0 {
		t.Fatal("old tables retained")
	}
}
