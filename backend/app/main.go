package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/BX-Research/ngadu-monorepo/backend/app/dal"
	"github.com/BX-Research/ngadu-monorepo/backend/app/database"
	"github.com/BX-Research/ngadu-monorepo/backend/app/routes"
	"github.com/BX-Research/ngadu-monorepo/backend/app/services"
	"github.com/BX-Research/ngadu-monorepo/backend/app/storage"
	"github.com/gofiber/fiber/v3"
)

func main() {
	databaseURL, err := database.URL()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := database.Open(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		ReadTimeout:   5 * time.Second, IdleTimeout: 60 * time.Second,
	})

	app.Get("/healthz", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/readyz", func(c fiber.Ctx) error {
		c.Set("Cache-Control", "no-store")
		ctx, cancel := context.WithTimeout(c.Context(), time.Second)
		defer cancel()
		if err := conn.Ping(ctx); err != nil {
			return c.Status(503).JSON(fiber.Map{"status": "unavailable"})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})
	var photos services.PhotoSigner
	config := storage.DevConfig()
	if err := config.ValidateLocal(); err != nil {
		log.Print("Profile photos unavailable: invalid local storage configuration")
	} else if store, err := storage.New(config); err != nil {
		log.Print("Profile photos unavailable: incomplete storage configuration")
	} else {
		photos = store
	}
	routes.LeaderboardRoutes(app, dal.NewLeaderboard(conn), photos)
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	routes.UserRoutes(app, dal.NewUsers(conn))
	log.Fatal(app.Listen("127.0.0.1:" + PORT))
}
