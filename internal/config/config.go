package config

import (
	"os"
	"time"
)

type Config struct {
	Server struct {
		Port string
	}
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	JWT struct {
		Secret string
		Expiry time.Duration
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func LoadConfig() *Config {
	cfg := &Config{}

	cfg.Server.Port = getEnv("PORT", "8080")

	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port = getEnv("DB_PORT", "3306")
	cfg.DB.User = getEnv("DB_USER", "huanlocal")
	cfg.DB.Password = getEnv("DB_PASSWORD", "pass123")
	cfg.DB.Name = getEnv("DB_NAME", "be")

	cfg.JWT.Secret = getEnv("JWT_SECRET", "indianman")
	cfg.JWT.Expiry = 12 * time.Hour

	return cfg
}
