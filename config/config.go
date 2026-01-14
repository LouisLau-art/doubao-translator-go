package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	APIKey         string
	APIURL         string
	Port           string
	GinMode        string
	CacheTTL       time.Duration
	CacheMaxSize   int
	MaxTextLength  int
	RateLimitRPM   int
	RateLimitBurst int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	cfg := &Config{
		APIKey:         getEnv("ARK_API_KEY", ""),
		APIURL:         getEnv("ARK_API_URL", "https://ark.cn-beijing.volces.com/api/v3/responses"),
		Port:           getEnv("PORT", "5000"),
		GinMode:        getEnv("GIN_MODE", "release"),
		CacheTTL:       getEnvAsDuration("CACHE_TTL", 3600*time.Second),
		CacheMaxSize:   getEnvAsInt("CACHE_MAX_SIZE", 1000),
		MaxTextLength:  getEnvAsInt("MAX_TEXT_LENGTH", 5000),
		RateLimitRPM:   getEnvAsInt("RATE_LIMIT_RPM", 30),
		RateLimitBurst: 30, // Default burst size
	}

	// Validate required configuration
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("ARK_API_KEY is required")
	}

	// Validate port
	if _, err := strconv.Atoi(cfg.Port); err != nil {
		return nil, fmt.Errorf("invalid PORT: %v", err)
	}

	return cfg, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// getEnvAsDuration gets an environment variable as duration
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	// Try to parse as seconds first
	if seconds, err := strconv.Atoi(valueStr); err == nil {
		return time.Duration(seconds) * time.Second
	}

	// Try to parse as duration string
	if duration, err := time.ParseDuration(valueStr); err == nil {
		return duration
	}

	return defaultValue
}