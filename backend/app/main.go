package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/BX-Research/ngadu-monorepo/backend/app/dal"
	"github.com/BX-Research/ngadu-monorepo/backend/app/database"
	"github.com/BX-Research/ngadu-monorepo/backend/app/routes"
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

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	routes.UserRoutes(app, dal.NewUsers(conn))
	log.Fatal(app.Listen(":" + PORT))
}
