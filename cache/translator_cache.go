package cache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// CacheItem represents a cached item with timestamp
type CacheItem struct {
	Value     string
	Timestamp time.Time
}

// TranslatorCache implements a thread-safe cache for translations
type TranslatorCache struct {
	store     sync.Map
	ttl       time.Duration
	maxSize   int
	size      int
	mu        sync.RWMutex
	cleanupCh chan struct{}
}

// NewTranslatorCache creates a new cache instance
func NewTranslatorCache(ttl time.Duration, maxSize int) *TranslatorCache {
	cache := &TranslatorCache{
		ttl:       ttl,
		maxSize:   maxSize,
		cleanupCh: make(chan struct{}),
	}

	// Start background cleanup goroutine
	go cache.startCleanup()

	return cache
}

// Get retrieves a value from cache
func (c *TranslatorCache) Get(key string) (string, bool) {
	if val, ok := c.store.Load(key); ok {
		item := val.(CacheItem)

		// Check if item has expired
		if time.Since(item.Timestamp) < c.ttl {
			return item.Value, true
		}

		// Remove expired item
		c.store.Delete(key)
		c.mu.Lock()
		c.size--
		c.mu.Unlock()
	}

	return "", false
}

// Set stores a value in cache
func (c *TranslatorCache) Set(key, value string) error {
	c.mu.RLock()
	if c.size >= c.maxSize {
		c.mu.RUnlock()
		return fmt.Errorf("cache is full (max size: %d)", c.maxSize)
	}
	c.mu.RUnlock()

	item := CacheItem{
		Value:     value,
		Timestamp: time.Now(),
	}

	c.store.Store(key, item)

	c.mu.Lock()
	c.size++
	c.mu.Unlock()

	return nil
}

// Delete removes a key from cache
func (c *TranslatorCache) Delete(key string) {
	c.store.Delete(key)

	c.mu.Lock()
	c.size--
	if c.size < 0 {
		c.size = 0
	}
	c.mu.Unlock()
}

// Clear removes all items from cache
func (c *TranslatorCache) Clear() {
	c.store.Range(func(key, value interface{}) bool {
		c.store.Delete(key)
		return true
	})

	c.mu.Lock()
	c.size = 0
	c.mu.Unlock()
}

// Size returns the current number of items in cache
func (c *TranslatorCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.size
}

// GetCacheKey generates a cache key from text and language codes
func GetCacheKey(text, source, target string) string {
	data := fmt.Sprintf("%s:%s:%s", source, target, text)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// startCleanup starts a background goroutine to clean expired items
func (c *TranslatorCache) startCleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanupExpired()
		case <-c.cleanupCh:
			return
		}
	}
}

// cleanupExpired removes all expired items from cache
func (c *TranslatorCache) cleanupExpired() {
	var expiredKeys []string

	c.store.Range(func(key, value interface{}) bool {
		item := value.(CacheItem)
		if time.Since(item.Timestamp) > c.ttl {
			expiredKeys = append(expiredKeys, key.(string))
		}
		return true
	})

	for _, key := range expiredKeys {
		c.Delete(key)
	}
}

// StopCleanup stops the background cleanup goroutine
func (c *TranslatorCache) StopCleanup() {
	close(c.cleanupCh)
}