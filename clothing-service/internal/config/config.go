package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppEnv     string
	HTTPPort   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	ClothingServiceURL string

	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket    string
	MinIOUseSSL    bool
	MinIOPublicURL string

	JWTSecret string
	JWTIssuer string
}

func Load() Config {
	return Config{
		AppEnv:   getEnv("APP_ENV", "development"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),

		// DB
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "app"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		ClothingServiceURL: getEnv("CLOTHING_SERVICE_URL", "http://localhost:8080"),

		// minIO
		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinIOBucket:    getEnv("MINIO_BUCKET", "clothes-images"),
		MinIOUseSSL:    getEnvAsBool("MINIO_USE_SSL", false),
		MinIOPublicURL: getEnv("MINIO_PUBLIC_URL", "http://localhost:9000"),

		JWTSecret: getEnv("JWT_SECRET", "super-secret-dev-key"),
		JWTIssuer: getEnv("JWT_ISSUER", "auth-service"),
	}
}

func (c Config) HTTPAddr() string {
	return ":" + c.HTTPPort
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value == "true"
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
