// Package storage signs private profile photos without making object-store reads.
package storage

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct{ Endpoint, Region, Bucket, AccessKey, SecretKey string }

func setting(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// DevConfig shares Compose's local defaults; it never consults ambient AWS credentials.
func DevConfig() Config {
	return Config{
		Endpoint:  "http://127.0.0.1:" + setting("DEV_S3_PORT", "8333"),
		Region:    setting("DEV_S3_REGION", "us-east-1"),
		Bucket:    setting("DEV_S3_BUCKET", "ngadu-dev"),
		AccessKey: setting("DEV_S3_ACCESS_KEY", "ngadu-local"),
		SecretKey: setting("DEV_S3_SECRET_KEY", "change-me-local-storage-only"),
	}
}

// ValidateLocal prevents development seed operations from targeting cloud storage.
func (c Config) ValidateLocal() error {
	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return fmt.Errorf("invalid local S3 endpoint")
	}
	ip := net.ParseIP(u.Hostname())
	if ip == nil || !ip.IsLoopback() {
		return fmt.Errorf("development storage requires a loopback endpoint")
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid development S3 port")
	}
	return nil
}

type Store struct {
	Client *s3.Client
	signer *s3.PresignClient
	bucket string
}

func New(c Config) (*Store, error) {
	u, err := url.Parse(c.Endpoint)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || u.RawQuery != "" || u.Fragment != "" || strings.Trim(u.Path, "/") != "" {
		return nil, fmt.Errorf("invalid S3 endpoint")
	}
	if c.Region == "" || c.Bucket == "" || c.AccessKey == "" || c.SecretKey == "" {
		return nil, fmt.Errorf("incomplete S3 configuration")
	}
	client := s3.New(s3.Options{
		Region: c.Region, BaseEndpoint: aws.String(c.Endpoint), UsePathStyle: true,
		Credentials: credentials.NewStaticCredentialsProvider(c.AccessKey, c.SecretKey, ""),
		HTTPClient:  &http.Client{Timeout: 10 * time.Second},
	})
	return &Store{Client: client, signer: s3.NewPresignClient(client), bucket: c.Bucket}, nil
}

func (s *Store) PhotoURL(ctx context.Context, key string) (string, error) {
	if strings.TrimSpace(key) == "" {
		return "", nil
	}
	result, err := s.signer.PresignGetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)}, func(o *s3.PresignOptions) { o.Expires = 15 * time.Minute })
	if err != nil {
		return "", err
	}
	return result.URL, nil
}
