// seed-storage uploads repository-owned fixtures only to loopback development S3.
package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/BX-Research/ngadu-monorepo/backend/app/avatars"
	"github.com/BX-Research/ngadu-monorepo/backend/app/config"
	"github.com/BX-Research/ngadu-monorepo/backend/app/database"
	"github.com/BX-Research/ngadu-monorepo/backend/app/storage"
	"github.com/BX-Research/ngadu-monorepo/backend/seeds"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/jackc/pgx/v5"
)

func missing(err error) bool {
	var api smithy.APIError
	return errors.As(err, &api) && (api.ErrorCode() == "NotFound" || api.ErrorCode() == "NoSuchBucket" || api.ErrorCode() == "NoSuchKey")
}

type photoWriter interface {
	HeadObject(context.Context, *s3.HeadObjectInput, ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	PutObject(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

// Conditional writes preserve concurrent edits. SeaweedFS can return 412 while
// its lock ring warms up, even for absent objects; recheck existence and retry
// within the seed's overall deadline rather than treating that as a stored file.
func seedPhoto(ctx context.Context, client photoWriter, bucket, key string, data []byte) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		_, err := client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})
		if err == nil {
			return nil
		}
		if !missing(err) {
			return fmt.Errorf("inspect local fixture: %w", err)
		}
		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket), Key: aws.String(key), Body: bytes.NewReader(data),
			ContentType: aws.String("image/jpeg"), IfNoneMatch: aws.String("*"),
			CacheControl: aws.String("private, no-cache"),
		})
		if err == nil {
			return nil
		}
		var api smithy.APIError
		if !errors.As(err, &api) || (api.ErrorCode() != "PreconditionFailed" && api.ErrorCode() != "ConditionalRequestConflict") {
			return fmt.Errorf("upload local fixture: %w", err)
		}
		timer := time.NewTimer(200 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}

type avatarUsers interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

func seedAvatars(ctx context.Context, db avatarUsers, store photoWriter, bucket string) error {
	fixtures, err := seeds.AvatarFixtures()
	if err != nil {
		return err
	}
	uids := make([]string, 0, len(fixtures))
	for _, fixture := range fixtures {
		uids = append(uids, fixture.FirebaseUID)
	}
	rows, err := db.Query(ctx, `SELECT id::text,firebase_uid FROM users WHERE firebase_uid=ANY($1)`, uids)
	if err != nil {
		return err
	}
	defer rows.Close()
	ids := map[string]string{}
	for rows.Next() {
		var id, uid string
		if err := rows.Scan(&id, &uid); err != nil {
			return err
		}
		ids[uid] = id
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(ids) != len(fixtures) {
		return fmt.Errorf("avatar users missing; run make db-seed first")
	}
	for _, fixture := range fixtures {
		key, err := avatars.Key(ids[fixture.FirebaseUID])
		if err != nil {
			return err
		}
		if err := seedPhoto(ctx, store, bucket, key, fixture.JPEG); err != nil {
			return err
		}
	}
	return nil
}

func run() error {
	if err := database.LoadEnv(); err != nil {
		return err
	}
	cfg := config.StorageFromEnv()
	if err := cfg.ValidateLocal(); err != nil {
		return err
	}
	store, err := storage.New(cfg)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err = store.Client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(cfg.Bucket)})
	if missing(err) {
		_, err = store.Client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(cfg.Bucket)})
	}
	if err != nil {
		return fmt.Errorf("prepare local bucket: %w", err)
	}
	databaseURL, err := database.LocalURL()
	if err != nil {
		return err
	}
	db, err := database.Open(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := seedAvatars(ctx, db, store.Client, cfg.Bucket); err != nil {
		return err
	}
	fmt.Println("Local profile photo fixtures ready (existing objects preserved).")
	return nil
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Storage seed failed:", err)
		os.Exit(1)
	}
}
