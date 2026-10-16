package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BX-Research/ngadu-monorepo/backend/internal/devdb"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

func newHandler(ping func(context.Context) error) *fiber.App {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		ReadTimeout:   5 * time.Second, IdleTimeout: 60 * time.Second,
	})

	app.Get("/healthz", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"}, fiber.MIMEApplicationJSON)
	})

	app.Get("/readyz", func(c fiber.Ctx) error {
		if err := ping(c.Context()); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "unavailable"}, fiber.MIMEApplicationJSON)
		}
		return c.JSON(fiber.Map{"status": "ready"}, fiber.MIMEApplicationJSON)
	})
	registerUserRoutes(app, findUser)
	app.All("/healthz", methodNotAllowed)
	app.All("/readyz", methodNotAllowed)
	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("404 page not found\n")
	})
	return app
}

func methodNotAllowed(c fiber.Ctx) error {
	c.Set("Allow", "GET, HEAD")
	return c.Status(fiber.StatusMethodNotAllowed).SendString("Method Not Allowed\n")
}

func pingDatabase(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	return conn.Ping(ctx)
}

func serve(ctx context.Context, listener net.Listener, app *fiber.App) error {
	return app.Listener(listener, fiber.ListenConfig{
		DisableStartupMessage: true,
		GracefulContext:       ctx,
		ShutdownTimeout:       5 * time.Second,
	})
}

func main() {
	connectionURL, err := devdb.URL()
	if err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("DATABASE_URL", connectionURL); err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Backend listening on http://%s", listener.Addr())
	if err := serve(ctx, listener, newHandler(pingDatabase)); err != nil {
		log.Fatal(err)
	}
}
