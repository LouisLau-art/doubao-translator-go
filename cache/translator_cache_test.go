package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCacheGetSet(t *testing.T) {
	// Create cache with 1 second TTL and max size 10
	c := NewTranslatorCache(1*time.Second, 10)

	// Set a value
	err := c.Set("test-key", "test-value")
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Get the value
	value, ok := c.Get("test-key")
	if !ok {
		t.Fatal("Expected to find key in cache")
	}

	if value != "test-value" {
		t.Errorf("Expected value to be 'test-value', got '%s'", value)
	}
}

func TestCacheExpiration(t *testing.T) {
	// Create cache with 100ms TTL
	c := NewTranslatorCache(100*time.Millisecond, 10)

	// Set value
	err := c.Set("test-key", "test-value")
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Try to get expired value
	value, ok := c.Get("test-key")
	if ok {
		t.Fatalf("Expected key to be expired, got value '%s'", value)
	}

	if value != "" {
		t.Errorf("Expected empty string for expired key, got '%s'", value)
	}
}

func TestCacheMaxSize(t *testing.T) {
	// Create cache with max size 2
	c := NewTranslatorCache(10*time.Second, 2)

	// Fill cache
	err := c.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Failed to set key1: %v", err)
	}

	err = c.Set("key2", "value2")
	if err != nil {
		t.Fatalf("Failed to set key2: %v", err)
	}

	// Try to set third key - should fail
	err = c.Set("key3", "value3")
	if err == nil {
		t.Fatal("Expected setting third key to fail with full cache")
	}

	// Check cache size
	if c.Size() != 2 {
		t.Errorf("Expected cache size to be 2, got %d", c.Size())
	}
}

func TestCacheClear(t *testing.T) {
	c := NewTranslatorCache(10*time.Second, 10)

	// Set some values
	err := c.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Failed to set key1: %v", err)
	}

	err = c.Set("key2", "value2")
	if err != nil {
		t.Fatalf("Failed to set key2: %v", err)
	}

	// Clear cache
	c.Clear()

	if c.Size() != 0 {
		t.Errorf("Expected cache size to be 0 after clear, got %d", c.Size())
	}

	// Try to get values
	_, ok := c.Get("key1")
	if ok {
		t.Error("Expected key1 to not exist after clear")
	}
}

func TestGetCacheKey(t *testing.T) {
	key := GetCacheKey("test text", "en", "zh")
	if key == "" {
		t.Fatal("Expected non-empty cache key")
	}

	// Check that same inputs produce same key
	key2 := GetCacheKey("test text", "en", "zh")
	if key != key2 {
		t.Errorf("Expected same key for same inputs, got different keys")
	}

	// Check that different inputs produce different keys
	key3 := GetCacheKey("different text", "en", "zh")
	if key == key3 {
		t.Errorf("Expected different keys for different inputs")
	}
}

func TestCacheDelete(t *testing.T) {
	c := NewTranslatorCache(10*time.Second, 10)

	// Set value
	err := c.Set("test-key", "test-value")
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Delete key
	c.Delete("test-key")

	// Try to get deleted key
	_, ok := c.Get("test-key")
	if ok {
		t.Error("Expected key to not exist after delete")
	}
}

func TestCacheConcurrentAccess(t *testing.T) {
	c := NewTranslatorCache(10*time.Second, 100)

	// Run concurrent writes and reads
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", i)
			value := fmt.Sprintf("value-%d", i)

			// Set value
			err := c.Set(key, value)
			if err != nil {
				t.Errorf("Failed to set key %s: %v", key, err)
				return
			}

			// Get value
			val, ok := c.Get(key)
			if !ok {
				t.Errorf("Failed to get key %s", key)
				return
			}

			if val != value {
				t.Errorf("Expected value %s for key %s, got %s", value, key, val)
			}
		}(i)
	}

	wg.Wait()

	// Verify final size
	if c.Size() != 10 {
		t.Errorf("Expected cache size to be 10, got %d", c.Size())
	}
}