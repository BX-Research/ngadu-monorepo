package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/BX-Research/ngadu-monorepo/backend/app/database/devdb"
	"github.com/BX-Research/ngadu-monorepo/backend/app/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	// dependency injection
	DATABASE_URL, err := devdb.URL()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := pgx.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

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
	routes.UserRoutes(app)
	app.Listen("127.0.0.1:" + PORT)
}
