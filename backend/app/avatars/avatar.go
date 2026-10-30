// Package avatars normalizes profile images and defines their storage convention.
package avatars

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	_ "golang.org/x/image/webp"
)

const (
	Size      = 512
	Quality   = 85
	MaxBytes  = 10 << 20
	MaxPixels = 25_000_000
)

var (
	ErrTooLarge    = errors.New("avatar exceeds 10 MiB")
	ErrDimensions  = errors.New("avatar exceeds 25 million pixels")
	ErrUnsupported = errors.New("avatar must be JPEG, PNG, GIF or WebP")
)

func Key(userID string) (string, error) {
	id, err := uuid.Parse(userID)
	if err != nil || id == uuid.Nil {
		return "", fmt.Errorf("invalid avatar user ID")
	}
	return "avatars/" + id.String() + ".jpg", nil
}

// Normalize validates dimensions before full decoding, applies JPEG EXIF
// orientation, center-crops, flattens alpha to white, and omits source metadata.
func Normalize(input io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(input, MaxBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read avatar: %w", err)
	}
	if len(data) > MaxBytes {
		return nil, ErrTooLarge
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, ErrUnsupported
	}
	switch format {
	case "jpeg", "png", "gif", "webp":
	default:
		return nil, ErrUnsupported
	}
	if cfg.Width <= 0 || cfg.Height <= 0 || int64(cfg.Width)*int64(cfg.Height) > MaxPixels {
		return nil, ErrDimensions
	}
	img, err := imaging.Decode(bytes.NewReader(data), imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("decode avatar: %w", err)
	}
	// Crop before resizing: Fill may allocate an enormous intermediate image
	// for valid, very narrow sources even when their pixel count is small.
	bounds := img.Bounds()
	side := min(bounds.Dx(), bounds.Dy())
	left, top := bounds.Min.X+(bounds.Dx()-side)/2, bounds.Min.Y+(bounds.Dy()-side)/2
	cropped := imaging.Crop(img, image.Rect(left, top, left+side, top+side))
	resized := imaging.Resize(cropped, Size, Size, imaging.Lanczos)
	flattened := image.NewRGBA(resized.Bounds())
	draw.Draw(flattened, flattened.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	draw.Draw(flattened, flattened.Bounds(), resized, resized.Bounds().Min, draw.Over)
	var output bytes.Buffer
	if err := jpeg.Encode(&output, flattened, &jpeg.Options{Quality: Quality}); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

type Writer interface {
	PutAvatar(context.Context, string, []byte) error
}
type Service struct{ writer Writer }

func New(writer Writer) *Service { return &Service{writer: writer} }

// Upload is intended for a future authenticated caller, which must authorize
// the user ID. It does not expose an HTTP endpoint or trust a client-supplied ID.
func (s *Service) Upload(ctx context.Context, userID string, input io.Reader) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, err := Key(userID); err != nil {
		return err
	}
	data, err := Normalize(input)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return s.writer.PutAvatar(ctx, userID, data)
}
