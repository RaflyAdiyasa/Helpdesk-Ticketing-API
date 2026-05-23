package config

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
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
	Minio struct {
		Endpoint   string
		AccessKey  string
		SecretKey  string
		BucketName string
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
	err := godotenv.Load()
	if err != nil {
		log.Warn(".env file not found, using environment variables or defaults")
	}
	cfg := &Config{}

	cfg.Server.Port = getEnv("PORT", "8080")

	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port = getEnv("DB_PORT", "3306")
	cfg.DB.User = getEnv("DB_USER", "helpdesk")
	cfg.DB.Password = getEnv("DB_PASSWORD", "helpdesk_password")
	cfg.DB.Name = getEnv("DB_NAME", "helpdesk")

	cfg.JWT.Secret = getEnv("JWT_SECRET", "indianman")
	cfg.JWT.Expiry = 12 * time.Hour

	cfg.Minio.Endpoint = getEnv("MINIO_ENDPOINT", "localhost:9000")
	cfg.Minio.AccessKey = getEnv("MINIO_ACCESSKEY", "minioadmin")
	cfg.Minio.SecretKey = getEnv("MINIO_SECRETKEY", "minioadmin123")
	cfg.Minio.BucketName = getEnv("MINIO_BUCKETNAME", "uploads")

	return cfg
}
