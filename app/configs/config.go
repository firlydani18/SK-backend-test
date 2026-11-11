package configs

import (
	"os"
	"strings"
)

type Config struct {
	Port     string
	DBPath   string
	SeedData bool
}

// Load reads environment configuration. A .env file is optional; environment
// variables supplied by the runtime take precedence.
func Load() Config {
	loadDotEnv(".env")
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./app.db"
	}
	seedData := strings.ToLower(os.Getenv("SEED_DATA")) != "false"
	return Config{Port: port, DBPath: dbPath, SeedData: seedData}
}

func loadDotEnv(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key, value = strings.TrimSpace(key), strings.Trim(strings.TrimSpace(value), "\"'")
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}
