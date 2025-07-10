package cache

import (
	"testing"
	"time"
)

func TestMemoryCache_BasicOperations(t *testing.T) {
	cache := NewMemoryCache[string](1 * time.Hour)
	defer cache.Stop()

	// Test Set and Get
	cache.Set("key1", "value1")
	value, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key1")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got '%v'", value)
	}

	// Test non-existent key
	_, found = cache.Get("nonexistent")
	if found {
		t.Error("Expected not to find nonexistent key")
	}

	// Test Delete
	cache.Delete("key1")
	_, found = cache.Get("key1")
	if found {
		t.Error("Expected key1 to be deleted")
	}
}

func TestMemoryCache_TTL(t *testing.T) {
	cache := NewMemoryCache[string](100 * time.Millisecond)
	defer cache.Stop()

	// Set item with short TTL
	cache.Set("expiring", "value")

	// Should be available immediately
	value, found := cache.Get("expiring")
	if !found || value != "value" {
		t.Error("Expected to find item immediately after setting")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired
	_, found = cache.Get("expiring")
	if found {
		t.Error("Expected item to be expired")
	}
}

func TestMemoryCache_CustomTTL(t *testing.T) {
	cache := NewMemoryCache[string](1 * time.Hour)
	defer cache.Stop()

	// Set item with custom short TTL
	cache.SetWithTTL("custom", "value", 100*time.Millisecond)

	// Should be available immediately
	value, found := cache.Get("custom")
	if !found || value != "value" {
		t.Error("Expected to find item immediately after setting")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired
	_, found = cache.Get("custom")
	if found {
		t.Error("Expected item to be expired")
	}
}

func TestMemoryCache_Size(t *testing.T) {
	cache := NewMemoryCache[string](1 * time.Hour)
	defer cache.Stop()

	// Initially empty
	if cache.Size() != 0 {
		t.Error("Expected cache to be empty initially")
	}

	// Add items
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	if cache.Size() != 2 {
		t.Errorf("Expected cache size to be 2, got %d", cache.Size())
	}

	// Clear cache
	cache.Clear()
	if cache.Size() != 0 {
		t.Error("Expected cache to be empty after clear")
	}
}

func TestMemoryCache_DifferentTypes(t *testing.T) {
	// Test with integer cache
	intCache := NewMemoryCache[int](1 * time.Hour)
	defer intCache.Stop()

	intCache.Set("number", 42)
	value, found := intCache.Get("number")
	if !found || value != 42 {
		t.Errorf("Expected 42, got %v", value)
	}

	// Test with struct cache
	type Person struct {
		Name string
		Age  int
	}

	personCache := NewMemoryCache[Person](1 * time.Hour)
	defer personCache.Stop()

	person := Person{Name: "Alice", Age: 30}
	personCache.Set("alice", person)
	retrievedPerson, found := personCache.Get("alice")
	if !found || retrievedPerson.Name != "Alice" || retrievedPerson.Age != 30 {
		t.Errorf("Expected %+v, got %+v", person, retrievedPerson)
	}
}
