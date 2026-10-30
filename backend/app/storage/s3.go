// Package storage signs private profile photos without making object-store reads.
package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/BX-Research/ngadu-monorepo/backend/app/avatars"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct{ EndpointHost, Endpoint, Region, Bucket, AccessKey, SecretKey string }

// ValidateLocal prevents development seed operations from targeting cloud storage.
func (c Config) ValidateLocal() error {
	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return fmt.Errorf("invalid local S3 endpoint")
	}
	ip := net.ParseIP(u.Hostname())
	if u.Hostname() != "localhost" && (ip == nil || !ip.IsLoopback()) {
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
	if c.EndpointHost != "" && !strings.EqualFold(c.EndpointHost, u.Hostname()) {
		return nil, fmt.Errorf("S3 endpoint host does not match URL")
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

func (s *Store) PhotoURL(ctx context.Context, userID string) (string, error) {
	key, err := avatars.Key(userID)
	if err != nil {
		return "", err
	}
	result, err := s.signer.PresignGetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)}, func(o *s3.PresignOptions) { o.Expires = 15 * time.Minute })
	if err != nil {
		return "", err
	}
	return result.URL, nil
}

// PutAvatar stores JPEG bytes supplied by the avatar normalization pipeline.
func (s *Store) PutAvatar(ctx context.Context, userID string, data []byte) error {
	key, err := avatars.Key(userID)
	if err != nil {
		return err
	}
	_, err = s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket), Key: aws.String(key), Body: bytes.NewReader(data),
		ContentType: aws.String("image/jpeg"), CacheControl: aws.String("private, no-cache"),
	})
	return err
}
