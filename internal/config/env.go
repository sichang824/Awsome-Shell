package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadEnv loads .env from current directory or project root.
func LoadEnv() {
	// Try current dir first
	_ = godotenv.Load(".env")
	// Try workspace root (where docker-compose usually is)
	cwd, _ := os.Getwd()
	for i := 0; i < 5; i++ {
		p := filepath.Join(cwd, ".env")
		if _, err := os.Stat(p); err == nil {
			_ = godotenv.Load(p)
			return
		}
		cwd = filepath.Dir(cwd)
		if cwd == "/" || cwd == "." {
			break
		}
	}
}

// GetEnv returns env var or default.
func GetEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
