package devenv

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func TestComposeAndGoReadSameCredentials(t *testing.T) {
	if os.Getenv("TEST_COMPOSE_DOTENV") == "" {
		t.Skip("set TEST_COMPOSE_DOTENV to check Docker Compose dotenv compatibility")
	}
	for _, prefix := range []string{"S3_", "DEV_S3_"} {
		t.Run(prefix, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, ".env")
			compose := filepath.Join(dir, "compose.yaml")
			if err := os.WriteFile(path, []byte(prefix+"SECRET_KEY='synthetic!secret $LITERAL'\n"), 0600); err != nil {
				t.Fatal(err)
			}
			if err := Configure(path, nil, "9000"); err != nil {
				t.Fatal(err)
			}
			values, err := godotenv.Read(path)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(compose, []byte("services:\n  storage:\n    image: alpine\n    environment:\n      TEST_SECRET: ${S3_SECRET_KEY}\n"), 0600); err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command("docker", "compose", "--env-file", path, "-f", compose, "config", "--environment")
			for _, value := range os.Environ() {
				if !strings.HasPrefix(value, "S3_SECRET_KEY=") {
					cmd.Env = append(cmd.Env, value)
				}
			}
			raw, err := cmd.Output()
			if err != nil {
				t.Fatal(err)
			}
			var composeSecret string
			for _, line := range strings.Split(string(raw), "\n") {
				if value, ok := strings.CutPrefix(line, "S3_SECRET_KEY="); ok {
					composeSecret = value
				}
			}
			if composeSecret != values["S3_SECRET_KEY"] || values["S3_SECRET_KEY"] != "synthetic!secret $LITERAL" {
				t.Fatal("Compose and Go interpreted fixture credentials differently")
			}

		})
	}
}
