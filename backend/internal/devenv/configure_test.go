package devenv

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigurePreservesSettingsAndWritesLocalEndpoint(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	original := "# keep this comment\nDATABASE_URL=postgres://unchanged\nAPI_TOKEN='keep $() literally'\nDEV_S3_BUCKET=edited-bucket\nDEV_S3_ACCESS_KEY=edited-key\nDEV_S3_SECRET_KEY=edited-secret\nS3_ENDPOINT_HOST=remote.invalid\nS3_ENDPOINT_URL=https://remote.invalid\n"
	if err := os.WriteFile(path, []byte(original), 0600); err != nil {
		t.Fatal(err)
	}
	defaults := map[string]string{"S3_BUCKET": "ngadu-dev", "S3_REGION": "us-east-1", "S3_ACCESS_KEY": "key", "S3_SECRET_KEY": "secret"}
	if err := Configure(path, defaults, "9000"); err != nil {
		t.Fatal(err)
	}
	first, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"# keep this comment", "DATABASE_URL=postgres://unchanged", "API_TOKEN='keep $() literally'", "S3_ENDPOINT_HOST=localhost", "S3_ENDPOINT_URL=http://localhost:9000", "S3_PORT=9000", "S3_BUCKET=edited-bucket", "S3_ACCESS_KEY=edited-key", "S3_SECRET_KEY=edited-secret"} {
		if !strings.Contains(string(first), text) {
			t.Fatalf("missing preserved/configured line %q", text)
		}
	}
	if err := Configure(path, defaults, "9000"); err != nil {
		t.Fatal(err)
	}
	second, _ := os.ReadFile(path)
	if string(first) != string(second) {
		t.Fatal("repeat setup changed the file")
	}
	if strings.Count(string(second), "S3_ENDPOINT_URL=") != 1 {
		t.Fatal("duplicate endpoint settings")
	}
	if err := Configure(path, defaults, "19000"); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "S3_ENDPOINT_URL=http://localhost:19000") {
		t.Fatal("isolated port not configured")
	}
}

func TestConfigureRejectsInvalidPortWithoutWriting(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("UNCHANGED=yes\n"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, port := range []string{"0", "65536", "bad"} {
		if err := Configure(path, nil, port); err == nil {
			t.Fatalf("accepted %q", port)
		}
	}
	data, _ := os.ReadFile(path)
	if string(data) != "UNCHANGED=yes\n" {
		t.Fatal("invalid setup wrote the file")
	}
}

func TestConfigureKeepsGenericSettingsOverLegacyDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	original := "S3_BUCKET=custom-bucket\nDEV_S3_BUCKET=old-bucket\nS3_SECRET_KEY='secret with spaces # and $()'\n"
	if err := os.WriteFile(path, []byte(original), 0600); err != nil {
		t.Fatal(err)
	}
	if err := Configure(path, map[string]string{"S3_BUCKET": "default", "S3_SECRET_KEY": "default"}, "9000"); err != nil {
		t.Fatal(err)
	}
	data, err := godotenv.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if data["S3_BUCKET"] != "custom-bucket" || data["S3_SECRET_KEY"] != "secret with spaces # and $()" {
		t.Fatal("generic settings were overwritten or incorrectly quoted")
	}
}

func TestConfigurePreservesCredentialSyntax(t *testing.T) {
	for _, prefix := range []string{"S3_", "DEV_S3_"} {
		t.Run(prefix, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), ".env")
			if err := os.WriteFile(path, []byte(prefix+"SECRET_KEY='synthetic!secret $LITERAL'\n"), 0600); err != nil {
				t.Fatal(err)
			}
			if err := Configure(path, nil, "9000"); err != nil {
				t.Fatal(err)
			}
			raw, _ := os.ReadFile(path)
			if !strings.Contains("\n"+string(raw), "\nS3_SECRET_KEY='synthetic!secret $LITERAL'") {
				t.Fatal("existing credential syntax was re-encoded")
			}
		})
	}
}

func TestConfigurePreservesMultilineValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	original := "DEV_S3_SECRET_KEY='first\nS3_PORT=1234\nlast'\nUNRELATED='before\nS3_ENDPOINT_URL=https://keep.invalid\nafter'\n"
	if err := os.WriteFile(path, []byte(original), 0600); err != nil {
		t.Fatal(err)
	}
	if err := Configure(path, nil, "9000"); err != nil {
		t.Fatal(err)
	}
	data, err := godotenv.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if data["S3_SECRET_KEY"] != "first\nS3_PORT=1234\nlast" || data["UNRELATED"] != "before\nS3_ENDPOINT_URL=https://keep.invalid\nafter" {
		t.Fatal("multiline value was changed")
	}
	if data["S3_PORT"] != "9000" || data["S3_ENDPOINT_URL"] != "http://localhost:9000" {
		t.Fatal("multiline contents masked real settings")
	}
}

func TestConfigureReplacesWholeMultilineEndpoint(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("S3_ENDPOINT_URL='https://old.invalid\nextra'\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := Configure(path, nil, "9000"); err != nil {
		t.Fatal(err)
	}
	data, err := godotenv.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if data["S3_ENDPOINT_URL"] != "http://localhost:9000" {
		t.Fatal("multiline endpoint not replaced")
	}
}
