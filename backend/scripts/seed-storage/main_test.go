package main

import (
	"context"
	"github.com/BX-Research/ngadu-monorepo/backend/app/database"
	"github.com/BX-Research/ngadu-monorepo/backend/internal/testdb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type fixtureStore struct {
	data      []byte
	transient bool
	puts      int
}

func (s *fixtureStore) HeadObject(context.Context, *s3.HeadObjectInput, ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	if s.data != nil {
		return &s3.HeadObjectOutput{}, nil
	}
	return nil, &smithy.GenericAPIError{Code: "NotFound"}
}
func (s *fixtureStore) PutObject(_ context.Context, input *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	s.puts++
	if s.transient {
		s.transient = false
		return nil, &smithy.GenericAPIError{Code: "PreconditionFailed"}
	}
	if s.data != nil {
		return nil, &smithy.GenericAPIError{Code: "PreconditionFailed"}
	}
	data := make([]byte, 32)
	n, _ := input.Body.Read(data)
	s.data = data[:n]
	return &s3.PutObjectOutput{}, nil
}
func TestSeedPhotoWaitsForWriteReadiness(t *testing.T) {
	store := &fixtureStore{transient: true}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := seedPhoto(ctx, store, "bucket", "photo.jpg", []byte("new image")); err != nil {
		t.Fatal(err)
	}
	if string(store.data) != "new image" || store.puts != 2 {
		t.Fatal("startup failure was not retried")
	}
}
func TestSeedPhotoPreservesExistingObject(t *testing.T) {
	store := &fixtureStore{data: []byte("edited image")}
	if err := seedPhoto(context.Background(), store, "bucket", "photo.jpg", []byte("fixture")); err != nil {
		t.Fatal(err)
	}
	if string(store.data) != "edited image" || store.puts != 0 {
		t.Fatal("existing object overwritten")
	}
}
func TestSeedPhotoStopsOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := seedPhoto(ctx, &fixtureStore{transient: true}, "bucket", "photo.jpg", []byte("fixture")); err == nil {
		t.Fatal("canceled seed kept running")
	}
}

type avatarObjects struct {
	objects     map[string][]byte
	contentType string
}

func (s *avatarObjects) HeadObject(_ context.Context, input *s3.HeadObjectInput, _ ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	if _, ok := s.objects[*input.Key]; ok {
		return &s3.HeadObjectOutput{}, nil
	}
	return nil, &smithy.GenericAPIError{Code: "NotFound"}
}
func (s *avatarObjects) PutObject(_ context.Context, input *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	data, err := io.ReadAll(input.Body)
	if err != nil {
		return nil, err
	}
	s.objects[*input.Key] = data
	s.contentType = *input.ContentType
	return &s3.PutObjectOutput{}, nil
}
func TestSeedAvatarsUsesStoredIDsAndPreservesObjects(t *testing.T) {
	db := testdb.Open(t)
	ctx := context.Background()
	root, err := database.Root()
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "backend/seeds/dev.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(ctx, string(data)); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(ctx, `DELETE FROM player_ratings WHERE user_id='00000000-0000-4000-8000-000000000201'; UPDATE users SET id='00000000-0000-4000-8000-000000000901' WHERE firebase_uid='seed-player-201'`); err != nil {
		t.Fatal(err)
	}
	editedKey := "avatars/00000000-0000-4000-8000-000000000901.jpg"
	store := &avatarObjects{objects: map[string][]byte{editedKey: []byte("edited")}}
	if err := seedAvatars(ctx, db, store, "bucket"); err != nil {
		t.Fatal(err)
	}
	if len(store.objects) != 14 || string(store.objects[editedKey]) != "edited" || store.contentType != "image/jpeg" {
		t.Fatal("seed avatar identity/preservation failed")
	}
	if _, ok := store.objects["avatars/00000000-0000-4000-8000-000000000201.jpg"]; ok {
		t.Fatal("used fixture UUID instead of actual stored ID")
	}
	if err := seedAvatars(ctx, db, store, "bucket"); err != nil {
		t.Fatal(err)
	}
	if len(store.objects) != 14 {
		t.Fatal("repeat seed duplicated objects")
	}
}
