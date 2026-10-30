package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
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
	if err := seedPhoto(ctx, store, "bucket", "photo.svg", []byte("new image")); err != nil {
		t.Fatal(err)
	}
	if string(store.data) != "new image" || store.puts != 2 {
		t.Fatal("startup failure was not retried")
	}
}
func TestSeedPhotoPreservesExistingObject(t *testing.T) {
	store := &fixtureStore{data: []byte("edited image")}
	if err := seedPhoto(context.Background(), store, "bucket", "photo.svg", []byte("fixture")); err != nil {
		t.Fatal(err)
	}
	if string(store.data) != "edited image" || store.puts != 0 {
		t.Fatal("existing object overwritten")
	}
}
func TestSeedPhotoStopsOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := seedPhoto(ctx, &fixtureStore{transient: true}, "bucket", "photo.svg", []byte("fixture")); err == nil {
		t.Fatal("canceled seed kept running")
	}
}
