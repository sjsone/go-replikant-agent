package cli

import (
	"bufio"
	"os"
	"strings"
)

// loadDotEnv reads a .env file (if present) and sets variables into the environment.
// Only variables that are not already set are overridden (env takes precedence).
func LoadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"`)
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

func EnvOrFlag(env, flag string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return flag
}
