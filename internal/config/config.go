package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port         string
	DatabaseHost string
	DatabasePort string
	DatabaseUser string
	DatabasePass string
	DatabaseName string
	DefaultCarID int
}

func Load() Config {
	loadDotEnv(".env")

	c := Config{
		Port:         envOr("PORT", "8080"),
		DatabaseHost: envOr("DATABASE_HOST", "teslamate-cnpg-rw"),
		DatabasePort: envOr("DATABASE_PORT", "5432"),
		DatabaseUser: envOr("DATABASE_USER", "teslamate"),
		DatabasePass: envOr("DATABASE_PASS", ""),
		DatabaseName: envOr("DATABASE_NAME", "teslamate"),
		DefaultCarID: 1,
	}
	if v := os.Getenv("DEFAULT_CAR_ID"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			c.DefaultCarID = id
		}
	}
	return c
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func loadDotEnv(path string) {
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
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
}
