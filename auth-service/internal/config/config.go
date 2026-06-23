package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
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

	AuthServiceURL string

	JWTSecret                string
	JWTIssuer                string
	JWTAccessTokenTTLMinutes int
}

func Load() Config {
	return Config{
		AppEnv:   getEnv("APP_ENV", "development"),
		HTTPPort: getEnv("HTTP_PORT", "8090"),

		// DB
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "app"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8090"),

		JWTSecret:                getEnv("JWT_SECRET", "super-secret-dev-key"),
		JWTIssuer:                getEnv("JWT_ISSUER", "auth-service"),
		JWTAccessTokenTTLMinutes: getEnvAsInt("JWT_ACCESS_TOKEN_TTL_MINUTES", 60),
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

func (c Config) JWTAccessTokenTTL() time.Duration {
	return time.Duration(c.JWTAccessTokenTTLMinutes) * time.Minute
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return n
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
