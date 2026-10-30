// Package devenv prepares local settings without rewriting unrelated environment entries.
package devenv

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var assignment = regexp.MustCompile(`^\s*(?:export\s+)?([A-Za-z_][A-Za-z0-9_]*)\s*=`)
var plainValue = regexp.MustCompile(`^[A-Za-z0-9_./:@-]*$`)

// closesQuote follows godotenv's quoted-value terminator rules.
func closesQuote(value string, quote byte) bool {
	for i := 0; i < len(value); i++ {
		if value[i] == quote && (i == 0 || value[i-1] != '\\') {
			return true
		}
	}
	return false
}

func Configure(path string, defaults map[string]string, port string) error {
	number, err := strconv.Atoi(port)
	if err != nil || number < 1 || number > 65535 {
		return fmt.Errorf("invalid S3_PORT")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	existing, err := godotenv.Unmarshal(string(raw))
	if err != nil {
		return fmt.Errorf("parse local environment: %w", err)
	}
	keys := []string{"S3_ENDPOINT_HOST", "S3_ENDPOINT_URL", "S3_PORT", "S3_BUCKET", "S3_REGION", "S3_ACCESS_KEY", "S3_SECRET_KEY"}
	values := map[string]string{"S3_ENDPOINT_HOST": "localhost", "S3_ENDPOINT_URL": "http://localhost:" + port, "S3_PORT": port}
	for _, key := range keys[3:] {
		value, ok := existing[key]
		if !ok {
			value, ok = existing["DEV_"+key]
		}
		if !ok {
			value = defaults[key]
		}
		values[key] = value
	}
	lines := []string{}
	seen := map[string]bool{}
	encode := func(key string) (string, error) {
		if plainValue.MatchString(values[key]) {
			return key + "=" + values[key], nil
		}
		return "", fmt.Errorf("default for %s must be a plain placeholder value", key)
	}
	var quote byte
	dropContinuation := false
	for _, line := range strings.Split(string(raw), "\n") {
		if quote != 0 {
			if !dropContinuation {
				lines = append(lines, line)
			}
			if closesQuote(line, quote) {
				quote = 0
				dropContinuation = false
			}
			continue
		}
		match := assignment.FindStringSubmatch(line)
		if len(match) > 1 {
			key := match[1]
			value := strings.TrimSpace(line[strings.IndexByte(line, '=')+1:])
			if len(value) > 0 && (value[0] == '\'' || value[0] == '"') && !closesQuote(value[1:], value[0]) {
				quote = value[0]
			}
			// Rename legacy credentials without re-encoding their quoted value. Both
			// Compose and Go must see exactly the same syntax, including punctuation.
			if strings.HasPrefix(key, "DEV_S3_") {
				generic := strings.TrimPrefix(key, "DEV_")
				if _, managed := values[generic]; managed {
					if _, present := existing[generic]; !present {
						at := strings.Index(line, key)
						line = line[:at] + generic + line[at+len(key):]
						key = generic
					}
				}
			}
			if _, ok := values[key]; ok {
				if key == "S3_ENDPOINT_HOST" || key == "S3_ENDPOINT_URL" || key == "S3_PORT" {
					dropContinuation = quote != 0
					if seen[key] {
						continue
					}
					line, err = encode(key)
					if err != nil {
						return err
					}
				}
				seen[key] = true
			}
		}
		lines = append(lines, line)
	}
	for _, key := range keys {
		if !seen[key] {
			line, err := encode(key)
			if err != nil {
				return err
			}
			lines = append(lines, line)
		}
	}
	result := strings.TrimRight(strings.Join(lines, "\n"), "\r\n") + "\n"
	if result == string(raw) {
		return nil
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	temp, err := os.CreateTemp(filepath.Dir(path), ".env.dev-*")
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())
	if err := temp.Chmod(info.Mode().Perm()); err != nil {
		temp.Close()
		return err
	}
	if _, err := temp.WriteString(result); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return os.Rename(temp.Name(), path)
}
