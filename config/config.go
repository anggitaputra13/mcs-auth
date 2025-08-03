package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	Host          string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	RedisURL      string
	JWTSecret     string
	JWTExpiration int
}

func LoadConfig() (*Config, error) {
	// Try to load .env file but don't fail if it doesn't exist
	_ = godotenv.Load()

	jwtExpiration, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:          getEnv("PORT", "8000"),
		Host:          getEnv("HOST", "localhost:8000"),
		DBHost:        getEnv("DB_HOST", "postgres"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "postgres"),
		DBName:        getEnv("DB_NAME", "auth_service"),
		DBSSLMode:     getEnv("DB_SSL_MODE", "disable"),
		RedisURL:      getEnv("REDIS_URL", "redis://redis:6379/0"),
		JWTSecret:     getEnv("JWT_SECRET", "your-very-secret-key"),
		JWTExpiration: jwtExpiration,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
