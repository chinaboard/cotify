package cache

import (
	"sync"
	"time"
)

// Item represents a cached item with expiration time
type Item struct {
	Value     interface{}
	ExpiresAt time.Time
}

// IsExpired checks if the cache item has expired
func (i *Item) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

// MemoryCache represents an in-memory cache with TTL support
type MemoryCache struct {
	cache       *sync.Map
	defaultTTL  time.Duration
	stopCleanup chan struct{}
	cleanupOnce sync.Once
}

// NewMemoryCache creates a new memory cache instance
func NewMemoryCache(defaultTTL time.Duration) *MemoryCache {
	cache := &MemoryCache{
		cache:       &sync.Map{},
		defaultTTL:  defaultTTL,
		stopCleanup: make(chan struct{}),
	}

	// Start cleanup goroutine
	go cache.startCleanup()

	return cache
}

// startCleanup periodically removes expired items from cache
func (c *MemoryCache) startCleanup() {
	ticker := time.NewTicker(1 * time.Hour) // Cleanup every hour
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanupExpired()
		case <-c.stopCleanup:
			return
		}
	}
}

// cleanupExpired removes all expired items from cache
func (c *MemoryCache) cleanupExpired() {
	now := time.Now()
	c.cache.Range(func(key, value interface{}) bool {
		if item, ok := value.(*Item); ok {
			if now.After(item.ExpiresAt) {
				c.cache.Delete(key)
			}
		}
		return true
	})
}

// Get retrieves an item from cache
// Returns the value and a boolean indicating if the item was found and not expired
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	if value, exists := c.cache.Load(key); exists {
		if item, ok := value.(*Item); ok {
			if !item.IsExpired() {
				return item.Value, true
			}
			// Item expired, remove it
			c.cache.Delete(key)
		}
	}
	return nil, false
}

// Set stores an item in cache with default TTL
func (c *MemoryCache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL stores an item in cache with custom TTL
func (c *MemoryCache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.cache.Store(key, &Item{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	})
}

// Delete removes an item from cache
func (c *MemoryCache) Delete(key string) {
	c.cache.Delete(key)
}

// Clear removes all items from cache
func (c *MemoryCache) Clear() {
	c.cache.Range(func(key, value interface{}) bool {
		c.cache.Delete(key)
		return true
	})
}

// Size returns the approximate number of items in cache
// Note: This is not guaranteed to be exact due to concurrent access
func (c *MemoryCache) Size() int {
	count := 0
	c.cache.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// Stop stops the cleanup goroutine
func (c *MemoryCache) Stop() {
	c.cleanupOnce.Do(func() {
		close(c.stopCleanup)
	})
}
