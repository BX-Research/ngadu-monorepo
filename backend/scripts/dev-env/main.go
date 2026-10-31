package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BX-Research/ngadu-monorepo/backend/app/database"
	"github.com/BX-Research/ngadu-monorepo/backend/internal/devenv"
	"github.com/joho/godotenv"
)

func run() error {
	root, err := database.Root()
	if err != nil {
		return err
	}
	defaults, err := godotenv.Read(filepath.Join(root, ".env.example"))
	if err != nil {
		return err
	}
	path := filepath.Join(root, ".env")
	existing, err := godotenv.Read(path)
	if err != nil {
		return err
	}
	port := os.Getenv("S3_PORT")
	if port == "" {
		port = existing["S3_PORT"]
	}
	if port == "" {
		port = "9000"
	}
	return devenv.Configure(path, defaults, port)
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Local environment setup failed:", err)
		os.Exit(1)
	}
}
