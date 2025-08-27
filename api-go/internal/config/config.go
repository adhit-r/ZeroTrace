package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Server configuration
	Port  int
	Host  string
	Debug bool

	// Database configuration
	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	// Redis configuration
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int

	// JWT configuration
	JWTSecret string
	JWTExpiry time.Duration

	// Rate limiting
	RateLimitRequests int
	RateLimitWindow   time.Duration

	// Logging
	LogLevel  string
	LogFormat string

	// Enrichment service
	EnrichmentServiceURL string
}

func Load() *Config {
	return &Config{
		// Server
		Port:  getEnvAsInt("API_PORT", 8080),
		Host:  getEnv("API_HOST", "0.0.0.0"),
		Debug: getEnvAsBool("API_MODE", "debug"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBName:     getEnv("DB_NAME", "zerotrace"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// Redis
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnvAsInt("REDIS_PORT", 6379),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "dev-secret-key-change-in-production"),
		JWTExpiry: getEnvAsDuration("JWT_EXPIRY", "24h"),

		// Rate limiting
		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvAsDuration("RATE_LIMIT_WINDOW", "1m"),

		// Logging
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),

		// Enrichment service
		EnrichmentServiceURL: getEnv("ENRICHMENT_SERVICE_URL", "http://localhost:8000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key, defaultValue string) bool {
	value := getEnv(key, defaultValue)
	return value == "true" || value == "debug"
}

func getEnvAsDuration(key, defaultValue string) time.Duration {
	value := getEnv(key, defaultValue)
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	return 24 * time.Hour // Default to 24 hours
}
