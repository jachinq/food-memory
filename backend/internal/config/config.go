package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	AppPort         int
	DBHost          string
	DBPort          int
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
	UploadDir       string
	PublicDir       string
	AccessToken     string
	MaxUploadBytes  int64
}

func Load() *Config {
	loadDotEnv()
	return &Config{
		AppPort:        envInt("APP_PORT", 8080),
		DBHost:         env("DB_HOST", "127.0.0.1"),
		DBPort:         envInt("DB_PORT", 5432),
		DBUser:         env("DB_USER", "shiyi"),
		DBPassword:     env("DB_PASSWORD", "shiyi_password"),
		DBName:         env("DB_NAME", "shiyi"),
		DBSSLMode:      env("DB_SSLMODE", "disable"),
		UploadDir:      env("UPLOAD_DIR", "./uploads"),
		PublicDir:      env("PUBLIC_DIR", "./public"),
		AccessToken:    env("APP_ACCESS_TOKEN", ""),
		MaxUploadBytes: 10 * 1024 * 1024,
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode,
	)
}

func loadDotEnv() {
	for _, path := range []string{".env", filepath.Join("backend", ".env")} {
		if applyDotEnvFile(path) {
			return
		}
	}
}

func applyDotEnvFile(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.Trim(strings.TrimSpace(val), `"'`)
		if key == "" {
			continue
		}
		if os.Getenv(key) == "" {
			_ = os.Setenv(key, val)
		}
	}
	return true
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
