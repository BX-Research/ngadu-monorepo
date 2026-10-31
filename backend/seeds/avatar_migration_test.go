package seeds

import (
	"bytes"
	"context"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BX-Research/ngadu-monorepo/backend/app/database"
	"github.com/BX-Research/ngadu-monorepo/backend/internal/testdb"
)

func TestDerivedAvatarMigrationPreservesProfilesAndRatings(t *testing.T) {
	const migration = "20261030000002_derived_avatars.sql"
	db := testdb.OpenBefore(t, migration)
	ctx := context.Background()
	if err := seedSQL(ctx, db); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(ctx, `UPDATE users SET display_name='Edited Name',profile_photo_key='legacy/photo.svg' WHERE firebase_uid='seed-player-201'; UPDATE player_ratings SET rating=2001 WHERE user_id='00000000-0000-4000-8000-000000000201'`); err != nil {
		t.Fatal(err)
	}
	root, err := database.Root()
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(root, "backend/migrations", migration))
	if err != nil {
		t.Fatal(err)
	}
	up, _, _ := strings.Cut(string(raw), "-- +goose Down")
	if _, err = db.Exec(ctx, up); err != nil {
		t.Fatal(err)
	}
	var count, rating int
	var name string
	if err = db.QueryRow(ctx, `SELECT count(*) FROM information_schema.columns WHERE table_schema=current_schema() AND table_name='users' AND column_name='profile_photo_key'`).Scan(&count); err != nil || count != 0 {
		t.Fatal("photo key column retained")
	}
	if err = db.QueryRow(ctx, `SELECT display_name,rating FROM users JOIN player_ratings ON users.id=user_id WHERE firebase_uid='seed-player-201'`).Scan(&name, &rating); err != nil || name != "Edited Name" || rating != 2001 {
		t.Fatal("migration changed profile or rating")
	}
}

func TestAvatarFixturesAreJPEG(t *testing.T) {
	fixtures, err := AvatarFixtures()
	if err != nil {
		t.Fatal(err)
	}
	if len(fixtures) != 14 {
		t.Fatalf("fixture count %d", len(fixtures))
	}
	for _, fixture := range fixtures {
		cfg, err := jpeg.DecodeConfig(bytes.NewReader(fixture.JPEG))
		if err != nil || cfg.Width != 512 || cfg.Height != 512 {
			t.Fatalf("invalid fixture %s", fixture.FirebaseUID)
		}
	}
}
