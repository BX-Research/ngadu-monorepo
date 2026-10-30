package avatars

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"
	"testing"
)

const userID = "00000000-0000-4000-8000-000000000201"

func pngImage(t *testing.T, img image.Image) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestKeyUsesCanonicalUserID(t *testing.T) {
	key, err := Key(strings.ToUpper(userID))
	if err != nil || key != "avatars/"+userID+".jpg" {
		t.Fatalf("key %q error %v", key, err)
	}
	for _, id := range []string{"", "../other", "not-a-uuid"} {
		if _, err := Key(id); err == nil {
			t.Fatalf("accepted %q", id)
		}
	}
}
func TestNormalizeJPEGAndTransparentPNG(t *testing.T) {
	transparent := image.NewNRGBA(image.Rect(0, 0, 64, 32))
	data, err := Normalize(bytes.NewReader(pngImage(t, transparent)))
	if err != nil {
		t.Fatal(err)
	}
	result, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatal("output is not JPEG:", err)
	}
	if result.Bounds().Dx() != 512 || result.Bounds().Dy() != 512 {
		t.Fatal("avatar dimensions changed")
	}
	r, g, b, _ := result.At(256, 256).RGBA()
	if r < 65000 || g < 65000 || b < 65000 {
		t.Fatal("transparent background is not white")
	}
	var source bytes.Buffer
	if err := jpeg.Encode(&source, result, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := Normalize(&source); err != nil {
		t.Fatal("JPEG input rejected:", err)
	}
}
func TestNormalizeOrientationAndStripsMetadata(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			c := color.RGBA{255, 0, 0, 255}
			if x >= 16 {
				c = color.RGBA{0, 0, 255, 255}
			}
			img.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100}); err != nil {
		t.Fatal(err)
	}
	// EXIF TIFF orientation=6: rotate 90 degrees clockwise.
	exif := []byte{'E', 'x', 'i', 'f', 0, 0, 'I', 'I', 42, 0, 8, 0, 0, 0, 1, 0, 0x12, 1, 3, 0, 1, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0}
	segment := []byte{0xff, 0xe1, 0, byte(len(exif) + 2)}
	segment = append(segment, exif...)
	encoded := buf.Bytes()
	encoded = append(append(append([]byte{}, encoded[:2]...), segment...), encoded[2:]...)
	data, err := Normalize(bytes.NewReader(encoded))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(data, []byte("Exif")) {
		t.Fatal("metadata retained")
	}
	result, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	r, _, b, _ := result.At(256, 80).RGBA()
	if r <= b {
		t.Fatal("EXIF orientation was not applied")
	}
	r, _, b, _ = result.At(256, 430).RGBA()
	if b <= r {
		t.Fatal("orientation changed pixel ordering")
	}
}
func TestNormalizeRejectsInvalidAndOversizedImages(t *testing.T) {
	for _, data := range [][]byte{[]byte("<svg></svg>"), []byte("not an image"), bytes.Repeat([]byte{'a'}, MaxBytes+1)} {
		if _, err := Normalize(bytes.NewReader(data)); err == nil {
			t.Fatal("invalid image accepted")
		}
	}
	// A valid PNG header declaring excessive pixel dimensions, without allocating them.
	data := pngImage(t, image.NewRGBA(image.Rect(0, 0, 1, 1)))
	binary.BigEndian.PutUint32(data[16:20], 100000)
	binary.BigEndian.PutUint32(data[20:24], 100000)
	binary.BigEndian.PutUint32(data[29:33], crc32.ChecksumIEEE(data[12:29]))
	if _, err := Normalize(bytes.NewReader(data)); !errors.Is(err, ErrDimensions) {
		t.Fatalf("large dimensions: %v", err)
	}
}

type writer struct {
	id   string
	data []byte
	err  error
}

func (w *writer) PutAvatar(_ context.Context, id string, data []byte) error {
	w.id = id
	w.data = data
	return w.err
}
func TestUploadNormalizesBeforeInjectedWrite(t *testing.T) {
	target := &writer{}
	service := New(target)
	input := pngImage(t, image.NewRGBA(image.Rect(0, 0, 12, 20)))
	if err := service.Upload(context.Background(), userID, bytes.NewReader(input)); err != nil {
		t.Fatal(err)
	}
	if target.id != userID {
		t.Fatal("writer received wrong user")
	}
	if _, err := jpeg.Decode(bytes.NewReader(target.data)); err != nil {
		t.Fatal("raw upload reached storage")
	}
	target.err = errors.New("storage offline")
	if err := service.Upload(context.Background(), userID, bytes.NewReader(input)); !errors.Is(err, target.err) {
		t.Fatal("storage failure lost")
	}
}
func TestUploadRejectsInvalidDataBeforeStorage(t *testing.T) {
	target := &writer{}
	service := New(target)
	if err := service.Upload(context.Background(), userID, strings.NewReader("bad")); err == nil {
		t.Fatal("invalid image accepted")
	}
	if target.data != nil {
		t.Fatal("invalid image reached storage")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := service.Upload(ctx, userID, io.NopCloser(strings.NewReader("bad"))); !errors.Is(err, context.Canceled) {
		t.Fatal("cancellation ignored")
	}
}

func TestKeyCanonicalizesUUIDLetters(t *testing.T) {
	key, err := Key("ABCDEFAB-1234-4ABC-8ABC-ABCDEFABCDEF")
	if err != nil || key != "avatars/abcdefab-1234-4abc-8abc-abcdefabcdef.jpg" {
		t.Fatal("UUID not canonicalized")
	}
}

func TestNormalizeGIFAndWebP(t *testing.T) {
	var gifData bytes.Buffer
	if err := gif.Encode(&gifData, image.NewRGBA(image.Rect(0, 0, 8, 8)), nil); err != nil {
		t.Fatal(err)
	}
	webp, err := os.ReadFile("testdata/source.webp")
	if err != nil {
		t.Fatal(err)
	}
	for _, data := range [][]byte{gifData.Bytes(), webp} {
		normalized, err := Normalize(bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := jpeg.Decode(bytes.NewReader(normalized)); err != nil {
			t.Fatal("output is not JPEG")
		}
	}
}

func TestNormalizeExtremeAspectRatios(t *testing.T) {
	for _, bounds := range []image.Rectangle{image.Rect(0, 0, 100000, 1), image.Rect(0, 0, 1, 100000)} {
		data := pngImage(t, image.NewNRGBA(bounds))
		normalized, err := Normalize(bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		cfg, err := jpeg.DecodeConfig(bytes.NewReader(normalized))
		if err != nil || cfg.Width != 512 || cfg.Height != 512 {
			t.Fatal("narrow avatar normalization failed")
		}
	}
}
