package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppEnv   string
	HTTPPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret string
	JWTIssuer string

	ClothingServiceBaseURL string
	MinIOPublicURL         string
}

func Load() Config {
	return Config{
		AppEnv:   getEnv("APP_ENV", "development"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "app"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret: getEnv("JWT_SECRET", "super-secret-dev-key"),
		JWTIssuer: getEnv("JWT_ISSUER", "auth-service"),

		ClothingServiceBaseURL: getEnv("CLOTHING_SERVICE_BASE_URL", "http://localhost:8080"),
		MinIOPublicURL:         getEnv("MINIO_PUBLIC_URL", "http://localhost:9000"),
	}
}

func (c Config) HTTPAddr() string {
	return ":" + c.HTTPPort
}

func (c Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
	)
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
