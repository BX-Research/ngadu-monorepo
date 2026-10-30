package storage_test

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/BX-Research/ngadu-monorepo/backend/app/avatars"
	"github.com/BX-Research/ngadu-monorepo/backend/app/config"
	"github.com/BX-Research/ngadu-monorepo/backend/app/database"
	"github.com/BX-Research/ngadu-monorepo/backend/app/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func testStore(t *testing.T) (*storage.Store, storage.Config) {
	t.Helper()
	endpoint := os.Getenv("TEST_S3_ENDPOINT")
	if endpoint == "" {
		t.Skip("set TEST_S3_ENDPOINT after seeding disposable local storage")
	}
	if err := database.LoadEnv(); err != nil {
		t.Fatal(err)
	}
	cfg := config.StorageFromEnv()
	cfg.Endpoint = endpoint
	parsed, err := url.Parse(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	cfg.EndpointHost = parsed.Hostname()
	if err := cfg.ValidateLocal(); err != nil {
		t.Fatal(err)
	}
	store, err := storage.New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return store, cfg
}

func TestPrivateFixturePhotos(t *testing.T) {
	store, cfg := testStore(t)
	client := &http.Client{Timeout: 5 * time.Second}
	for _, tc := range []struct {
		name, id string
		signed   bool
		status   int
	}{
		{"signed JPEG", "00000000-0000-4000-8000-000000000201", true, 200},
		{"private JPEG", "00000000-0000-4000-8000-000000000201", false, 403},
		{"missing avatar", "00000000-0000-4000-8000-000000000205", true, 404},
	} {
		t.Run(tc.name, func(t *testing.T) {
			key, err := avatars.Key(tc.id)
			if err != nil {
				t.Fatal(err)
			}
			raw := cfg.Endpoint + "/" + cfg.Bucket + "/" + key
			if tc.signed {
				raw, err = store.PhotoURL(context.Background(), tc.id)
				if err != nil {
					t.Fatal(err)
				}
			}
			response, err := client.Get(raw)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if response.StatusCode != tc.status {
				t.Fatalf("status %d want %d", response.StatusCode, tc.status)
			}
			if tc.status == 200 {
				if response.Header.Get("Content-Type") != "image/jpeg" {
					t.Fatal("wrong content type")
				}
				image, err := jpeg.DecodeConfig(response.Body)
				if err != nil || image.Width != 512 || image.Height != 512 {
					t.Fatal("fixture is not a normalized JPEG")
				}
			}
		})
	}
}

func TestAvatarPipelineUploadsAndReplacesJPEG(t *testing.T) {
	store, cfg := testStore(t)
	ctx := context.Background()
	id := uuid.NewString()
	key, _ := avatars.Key(id)
	defer store.Client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(cfg.Bucket), Key: aws.String(key)})
	uploader := avatars.New(store)
	input := image.NewNRGBA(image.Rect(0, 0, 30, 60))
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, input); err != nil {
		t.Fatal(err)
	}
	if err := uploader.Upload(ctx, id, &encoded); err != nil {
		t.Fatal(err)
	}
	raw, err := store.PhotoURL(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	get := func() []byte {
		response, err := client.Get(raw)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()
		if response.StatusCode != 200 || response.Header.Get("Content-Type") != "image/jpeg" || !strings.Contains(response.Header.Get("Cache-Control"), "no-cache") {
			t.Fatal("avatar object headers incorrect")
		}
		data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		cfg, err := jpeg.DecodeConfig(bytes.NewReader(data))
		if err != nil || cfg.Width != 512 || cfg.Height != 512 {
			t.Fatal("upload was not normalized")
		}
		return data
	}
	first := get()
	draw.Draw(input, input.Bounds(), &image.Uniform{C: color.Black}, image.Point{}, draw.Src)
	if err := png.Encode(&encoded, input); err != nil {
		t.Fatal(err)
	}
	if err := uploader.Upload(ctx, id, &encoded); err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(first, get()) {
		t.Fatal("replacing avatar did not replace the existing key")
	}
}
