package cache

import (
	"sync"
	"time"
)

// Item represents a cached item with expiration time
type Item[T any] struct {
	Value     T
	ExpiresAt time.Time
}

// IsExpired checks if the cache item has expired
func (i *Item[T]) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

// MemoryCache represents an in-memory cache with TTL support
type MemoryCache[T any] struct {
	cache       *sync.Map
	defaultTTL  time.Duration
	stopCleanup chan struct{}
	cleanupOnce sync.Once
}

// NewMemoryCache creates a new memory cache instance
func NewMemoryCache[T any](defaultTTL time.Duration) *MemoryCache[T] {
	cache := &MemoryCache[T]{
		cache:       &sync.Map{},
		defaultTTL:  defaultTTL,
		stopCleanup: make(chan struct{}),
	}

	// Start cleanup goroutine
	go cache.startCleanup()

	return cache
}

// startCleanup periodically removes expired items from cache
func (c *MemoryCache[T]) startCleanup() {
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
func (c *MemoryCache[T]) cleanupExpired() {
	now := time.Now()
	c.cache.Range(func(key, value interface{}) bool {
		if item, ok := value.(*Item[T]); ok {
			if now.After(item.ExpiresAt) {
				c.cache.Delete(key)
			}
		}
		return true
	})
}

// Get retrieves an item from cache
// Returns the value and a boolean indicating if the item was found and not expired
func (c *MemoryCache[T]) Get(key string) (T, bool) {
	var zero T
	if value, exists := c.cache.Load(key); exists {
		if item, ok := value.(*Item[T]); ok {
			if !item.IsExpired() {
				return item.Value, true
			}
			// Item expired, remove it
			c.cache.Delete(key)
		}
	}
	return zero, false
}

// Set stores an item in cache with default TTL
func (c *MemoryCache[T]) Set(key string, value T) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL stores an item in cache with custom TTL
func (c *MemoryCache[T]) SetWithTTL(key string, value T, ttl time.Duration) {
	c.cache.Store(key, &Item[T]{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	})
}

// Delete removes an item from cache
func (c *MemoryCache[T]) Delete(key string) {
	c.cache.Delete(key)
}

// Clear removes all items from cache
func (c *MemoryCache[T]) Clear() {
	c.cache.Range(func(key, value interface{}) bool {
		c.cache.Delete(key)
		return true
	})
}

// Size returns the approximate number of items in cache
// Note: This is not guaranteed to be exact due to concurrent access
func (c *MemoryCache[T]) Size() int {
	count := 0
	c.cache.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// Stop stops the cleanup goroutine
func (c *MemoryCache[T]) Stop() {
	c.cleanupOnce.Do(func() {
		close(c.stopCleanup)
	})
}
