package config

import "testing"

func TestLoadUsesGenericS3Settings(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example.invalid/ngadu")
	t.Setenv("PORT", "8181")
	t.Setenv("S3_ENDPOINT_HOST", "storage.example")
	t.Setenv("S3_ENDPOINT_URL", "https://storage.example:9443")
	t.Setenv("S3_BUCKET", "profiles")
	t.Setenv("S3_REGION", "test-region")
	t.Setenv("S3_ACCESS_KEY", "key")
	t.Setenv("S3_SECRET_KEY", "secret")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Address != "127.0.0.1:8181" || cfg.S3.Endpoint != "https://storage.example:9443" || cfg.S3.EndpointHost != "storage.example" || cfg.S3.Bucket != "profiles" || cfg.S3.Region != "test-region" || cfg.S3.AccessKey != "key" || cfg.S3.SecretKey != "secret" {
		t.Fatal("configuration did not use generic environment settings")
	}
}
