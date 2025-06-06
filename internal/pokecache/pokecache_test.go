package pokecache

import (
	"testing"
	"time"
)

func TestCacheAdd(t *testing.T) {
	cache := NewCache(time.Minute)

	key := "test-key"
	val := []byte("test-value")

	cache.Add(key, val)

	// Check if data was added
	cache.mu.Lock()
	entry, exists := cache.data[key]
	cache.mu.Unlock()

	if !exists {
		t.Error("Expected key to exist in cache")
	}

	if string(entry.val) != string(val) {
		t.Errorf("Expected value %s, got %s", string(val), string(entry.val))
	}

	if entry.createdAt.IsZero() {
		t.Error("Expected createdAt to be set")
	}
}

func TestCacheGet(t *testing.T) {
	cache := NewCache(time.Minute)

	key := "test-key"
	val := []byte("test-value")

	// Test getting non-existent key
	_, exists := cache.Get("non-existent")
	if exists {
		t.Error("Expected false for non-existent key")
	}

	// Add data and test getting existing key
	cache.Add(key, val)

	retrieved, exists := cache.Get(key)
	if !exists {
		t.Error("Expected true for existing key")
	}

	if string(retrieved) != string(val) {
		t.Errorf("Expected value %s, got %s", string(val), string(retrieved))
	}
}

func TestCacheGetEmpty(t *testing.T) {
	cache := NewCache(time.Minute)

	val, exists := cache.Get("any-key")
	if exists {
		t.Error("Expected false for empty cache")
	}

	if val != nil {
		t.Error("Expected nil value for non-existent key")
	}
}

func TestCacheAddMultiple(t *testing.T) {
	cache := NewCache(time.Minute)

	testData := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
		"key3": []byte("value3"),
	}

	// Add multiple entries
	for k, v := range testData {
		cache.Add(k, v)
	}

	// Verify all entries exist
	for k, expectedVal := range testData {
		val, exists := cache.Get(k)
		if !exists {
			t.Errorf("Expected key %s to exist", k)
		}
		if string(val) != string(expectedVal) {
			t.Errorf("Expected value %s for key %s, got %s", string(expectedVal), k, string(val))
		}
	}
}

func TestCacheOverwrite(t *testing.T) {
	cache := NewCache(time.Minute)

	key := "test-key"
	val1 := []byte("value1")
	val2 := []byte("value2")

	// Add first value
	cache.Add(key, val1)
	retrieved1, _ := cache.Get(key)

	if string(retrieved1) != string(val1) {
		t.Errorf("Expected first value %s, got %s", string(val1), string(retrieved1))
	}

	// Overwrite with second value
	cache.Add(key, val2)
	retrieved2, _ := cache.Get(key)

	if string(retrieved2) != string(val2) {
		t.Errorf("Expected second value %s, got %s", string(val2), string(retrieved2))
	}
}

func TestReapLoop(t *testing.T) {
	// Use a very short interval for testing
	interval := 50 * time.Millisecond
	cache := NewCache(interval)

	key := "test-key"
	val := []byte("test-value")

	// Add data
	cache.Add(key, val)

	// Verify data exists
	_, exists := cache.Get(key)
	if !exists {
		t.Error("Expected data to exist immediately after adding")
	}

	// Wait for reap loop to run (should be longer than interval)
	time.Sleep(interval + 20*time.Millisecond)

	// Data should be expired and removed
	_, exists = cache.Get(key)
	if exists {
		t.Error("Expected data to be reaped after expiration")
	}
}

func TestReapLoopKeepsRecentData(t *testing.T) {
	interval := 100 * time.Millisecond
	cache := NewCache(interval)

	key := "test-key"
	val := []byte("test-value")

	// Add data
	cache.Add(key, val)

	// Wait less than the interval
	time.Sleep(interval / 2)

	// Data should still exist
	_, exists := cache.Get(key)
	if !exists {
		t.Error("Expected recent data to still exist")
	}
}

func TestReapLoopSelectiveExpiration(t *testing.T) {
	interval := 60 * time.Millisecond
	cache := NewCache(interval)

	// Add first entry
	cache.Add("old-key", []byte("old-value"))

	// Wait half the interval
	time.Sleep(interval / 2)

	// Add second entry
	cache.Add("new-key", []byte("new-value"))

	// Wait for the first entry to expire but not the second
	time.Sleep(interval)

	// Old entry should be expired
	_, oldExists := cache.Get("old-key")
	if oldExists {
		t.Error("Expected old entry to be expired")
	}

	// New entry should still exist
	_, newExists := cache.Get("new-key")
	if !newExists {
		t.Error("Expected new entry to still exist")
	}
}

func TestCacheConcurrency(t *testing.T) {
	cache := NewCache(time.Minute)

	// Test concurrent adds and gets
	done := make(chan bool)

	// Goroutine 1: Add data
	go func() {
		for i := 0; i < 100; i++ {
			cache.Add("key", []byte("value"))
		}
		done <- true
	}()

	// Goroutine 2: Get data
	go func() {
		for i := 0; i < 100; i++ {
			cache.Get("key")
		}
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	// If we get here without panicking, the test passes
}
