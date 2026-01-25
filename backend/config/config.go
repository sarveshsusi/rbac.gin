package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
	FrontendURL string     // ✅ ADD
	Mail        MailConfig // ✅ ADD
}

type ServerConfig struct {
	Port string
	Env  string // development | production
}

type DatabaseConfig struct {
	URL string
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

/* =====================
   Mail (Mailtrap)
===================== */

type MailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		JWT: JWTConfig{
			AccessSecret:  getEnv("JWT_ACCESS_SECRET", "access-secret"),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", "refresh-secret"),
			AccessExpiry:  15 * time.Minute,
			RefreshExpiry: 7 * 24 * time.Hour,
		},

		// ✅ FRONTEND URL (for reset links)
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),

		// ✅ MAILTRAP CONFIG
		Mail: MailConfig{
			Host:     getEnv("MAIL_HOST", ""),
			Port:     getEnvAsInt("MAIL_PORT", 587),
			Username: getEnv("MAIL_USERNAME", ""),
			Password: getEnv("MAIL_PASSWORD", ""),
			From:     getEnv("MAIL_FROM", "rbac@app.com"),
		},
	}
}

/* =====================
   Helpers
===================== */

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
