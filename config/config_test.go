package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Test with environment variables
	origAPIKey := os.Getenv("ARK_API_KEY")
	origAPIURL := os.Getenv("ARK_API_URL")
	origPort := os.Getenv("PORT")

	// Set test environment variables
	os.Setenv("ARK_API_KEY", "test-key-123")
	os.Setenv("PORT", "8080")

	// Load config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify values
	if cfg.APIKey != "test-key-123" {
		t.Errorf("Expected APIKey to be 'test-key-123', got '%s'", cfg.APIKey)
	}

	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got '%s'", cfg.Port)
	}

	if cfg.APIURL != "https://ark.cn-beijing.volces.com/api/v3/responses" {
		t.Errorf("Expected default APIURL, got '%s'", cfg.APIURL)
	}

	if cfg.GinMode != "release" {
		t.Errorf("Expected default GinMode to be 'release', got '%s'", cfg.GinMode)
	}

	if cfg.CacheTTL != 3600*time.Second {
		t.Errorf("Expected default CacheTTL to be 3600 seconds, got %v", cfg.CacheTTL)
	}

	if cfg.CacheMaxSize != 1000 {
		t.Errorf("Expected default CacheMaxSize to be 1000, got %d", cfg.CacheMaxSize)
	}

	if cfg.MaxTextLength != 5000 {
		t.Errorf("Expected default MaxTextLength to be 5000, got %d", cfg.MaxTextLength)
	}

	if cfg.RateLimitRPM != 30 {
		t.Errorf("Expected default RateLimitRPM to be 30, got %d", cfg.RateLimitRPM)
	}

	// Cleanup
	if origAPIKey != "" {
		os.Setenv("ARK_API_KEY", origAPIKey)
	} else {
		os.Unsetenv("ARK_API_KEY")
	}

	if origAPIURL != "" {
		os.Setenv("ARK_API_URL", origAPIURL)
	} else {
		os.Unsetenv("ARK_API_URL")
	}

	if origPort != "" {
		os.Setenv("PORT", origPort)
	} else {
		os.Unsetenv("PORT")
	}
}

func TestLoadConfigMissingAPIKey(t *testing.T) {
	// Save original value
	origAPIKey := os.Getenv("ARK_API_KEY")
	defer func() {
		if origAPIKey != "" {
			os.Setenv("ARK_API_KEY", origAPIKey)
		} else {
			os.Unsetenv("ARK_API_KEY")
		}
	}()

	// Clear API key
	os.Unsetenv("ARK_API_KEY")

	// Should fail to load config
	cfg, err := Load()
	if err == nil {
		t.Fatal("Expected error when API key is missing, got nil")
	}

	if cfg != nil {
		t.Fatal("Expected nil config when API key is missing")
	}
}

func TestLoadConfigWithCustomAPIURL(t *testing.T) {
	// Save original value
	origAPIKey := os.Getenv("ARK_API_KEY")
	origAPIURL := os.Getenv("ARK_API_URL")
	defer func() {
		if origAPIKey != "" {
			os.Setenv("ARK_API_KEY", origAPIKey)
		} else {
			os.Unsetenv("ARK_API_KEY")
		}

		if origAPIURL != "" {
			os.Setenv("ARK_API_URL", origAPIURL)
		} else {
			os.Unsetenv("ARK_API_URL")
		}
	}()

	// Set test values
	os.Setenv("ARK_API_KEY", "test-key-123")
	os.Setenv("ARK_API_URL", "https://custom.api.url/v1/translate")

	// Load config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.APIURL != "https://custom.api.url/v1/translate" {
		t.Errorf("Expected APIURL to be 'https://custom.api.url/v1/translate', got '%s'", cfg.APIURL)
	}
}

func TestLoadConfigWithCustomCacheConfig(t *testing.T) {
	// Save original values
	origAPIKey := os.Getenv("ARK_API_KEY")
	origCacheTTL := os.Getenv("CACHE_TTL")
	origCacheMaxSize := os.Getenv("CACHE_MAX_SIZE")
	defer func() {
		if origAPIKey != "" {
			os.Setenv("ARK_API_KEY", origAPIKey)
		} else {
			os.Unsetenv("ARK_API_KEY")
		}

		if origCacheTTL != "" {
			os.Setenv("CACHE_TTL", origCacheTTL)
		} else {
			os.Unsetenv("CACHE_TTL")
		}

		if origCacheMaxSize != "" {
			os.Setenv("CACHE_MAX_SIZE", origCacheMaxSize)
		} else {
			os.Unsetenv("CACHE_MAX_SIZE")
		}
	}()

	// Set test values
	os.Setenv("ARK_API_KEY", "test-key-123")
	os.Setenv("CACHE_TTL", "7200")  // 2 hours in seconds
	os.Setenv("CACHE_MAX_SIZE", "5000")

	// Load config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.CacheTTL != 7200*time.Second {
		t.Errorf("Expected CacheTTL to be 7200 seconds, got %v", cfg.CacheTTL)
	}

	if cfg.CacheMaxSize != 5000 {
		t.Errorf("Expected CacheMaxSize to be 5000, got %d", cfg.CacheMaxSize)
	}
}

func TestGetEnvWithDefault(t *testing.T) {
	// Save original value
	origTestVar := os.Getenv("TEST_VAR")
	defer func() {
		if origTestVar != "" {
			os.Setenv("TEST_VAR", origTestVar)
		} else {
			os.Unsetenv("TEST_VAR")
		}
	}()

	// Unset variable to test default
	os.Unsetenv("TEST_VAR")

	// Should return default
	val := getEnv("TEST_VAR", "default-value")
	if val != "default-value" {
		t.Errorf("Expected 'default-value', got '%s'", val)
	}

	// Set variable to test actual value
	os.Setenv("TEST_VAR", "actual-value")
	val = getEnv("TEST_VAR", "default-value")
	if val != "actual-value" {
		t.Errorf("Expected 'actual-value', got '%s'", val)
	}
}