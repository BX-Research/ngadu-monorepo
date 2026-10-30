// Package config loads application configuration independently of service wiring.
package config

import (
	"net"
	"os"

	"github.com/BX-Research/ngadu-monorepo/backend/app/database"
	"github.com/BX-Research/ngadu-monorepo/backend/app/storage"
)

type Config struct {
	DatabaseURL string
	Address     string
	S3          storage.Config
}

func Load() (Config, error) {
	databaseURL, err := database.URL()
	if err != nil {
		return Config{}, err
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Config{DatabaseURL: databaseURL, Address: net.JoinHostPort("127.0.0.1", port), S3: StorageFromEnv()}, nil
}

// StorageFromEnv reads explicit S3 settings; callers load their environment first.
// Local defaults belong to development setup, not the application client.
func StorageFromEnv() storage.Config {
	return storage.Config{
		EndpointHost: os.Getenv("S3_ENDPOINT_HOST"), Endpoint: os.Getenv("S3_ENDPOINT_URL"),
		Region: os.Getenv("S3_REGION"), Bucket: os.Getenv("S3_BUCKET"),
		AccessKey: os.Getenv("S3_ACCESS_KEY"), SecretKey: os.Getenv("S3_SECRET_KEY"),
	}
}
