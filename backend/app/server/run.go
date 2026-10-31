package server

import (
	"context"
	"log"

	"github.com/BX-Research/ngadu-monorepo/backend/app/config"
	"github.com/BX-Research/ngadu-monorepo/backend/app/dal"
	"github.com/BX-Research/ngadu-monorepo/backend/app/database"
	"github.com/BX-Research/ngadu-monorepo/backend/app/storage"
)

// Run is the composition root: it owns configuration, concrete clients and cleanup.
func Run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	db, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	deps := Dependencies{Users: dal.NewUsers(db), Leaderboard: dal.NewLeaderboard(db), Ready: db.Ping}
	if photos, err := storage.New(cfg.S3); err != nil {
		log.Print("Profile photos unavailable: invalid S3 configuration")
	} else {
		deps.Photos = photos
	}
	return New(deps).Listen(cfg.Address)
}
