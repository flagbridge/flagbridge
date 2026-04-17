package config

import (
	"os"
	"strings"
)

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	APIKeySalt     string
	AllowedOrigins []string
	SentryDSN      string
	EncryptionKey  string
}

func Load() *Config {
	return &Config{
		Port:           envOr("PORT", "8080"),
		DatabaseURL:    envOr("DATABASE_URL", ""),
		JWTSecret:      envOr("JWT_SECRET", ""),
		APIKeySalt:     envOr("API_KEY_SALT", ""),
		AllowedOrigins: parseOrigins(envOr("ALLOWED_ORIGINS", "https://vozes.social,https://admin.flagbridge.io,http://localhost:3000")),
		SentryDSN:      envOr("SENTRY_DSN", ""),
		EncryptionKey:  envOr("ENCRYPTION_KEY", ""),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseOrigins(s string) []string {
	parts := strings.Split(s, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			origins = append(origins, t)
		}
	}
	return origins
}
