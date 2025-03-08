// Package pokecache provides a simple, thread-safe in-memory caching system with
// automatic expiration of entries based on a configurable time interval.
// It is used primarily to store API responses to reduce the number of external
// API calls and improve application performance.
package pokecache

import (
	"sync"
	"time"
)

// Cache represents an in-memory key-value store with automatic entry expiration.
// It uses a mutex to ensure thread-safety for concurrent operations, making it
// suitable for use in concurrent applications.
type Cache struct {
	cache map[string]cacheEntry // Internal map storing cached data
	mu    sync.Mutex            // Mutex for thread-safe operations
}

// cacheEntry represents a single item in the cache.
// Each entry contains the cached value as a byte slice and a timestamp
// indicating when it was created, used to determine expiration.
type cacheEntry struct {
	val       []byte    // The cached data as a byte slice
	createdAt time.Time // Timestamp when the entry was created
}

// NewCache creates and initializes a new Cache with automatic cleanup.
// It starts a background goroutine that periodically removes expired entries
// based on the provided interval duration.
//
// Parameters:
//   - interval: The time duration after which cache entries are considered expired
//
// Returns:
//   - A pointer to the newly created Cache
func NewCache(interval time.Duration) *Cache {
	c := &Cache{cache: make(map[string]cacheEntry)}
	go c.reapLoop(interval)
	return c
}

// Add stores a value in the cache with the specified key.
// If the key already exists, its value will be overwritten with the new value.
// The entry is timestamped with the current UTC time for expiration tracking.
//
// Parameters:
//   - key: The string key to associate with the value
//   - val: The byte slice value to store in the cache
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheEntry{
		val:       val,
		createdAt: time.Now().UTC(),
	}
}

// Get retrieves a value from the cache by its key.
// It returns the value and a boolean indicating whether the key was found.
//
// Parameters:
//   - key: The string key to look up
//
// Returns:
//   - The byte slice value associated with the key
//   - A boolean indicating whether the key was found in the cache (true if found)
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.cache[key]
	return value.val, ok
}

// reapLoop runs in a separate goroutine and periodically triggers
// the cleanup of expired cache entries.
//
// Parameters:
//   - interval: The time duration between cleanup cycles
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(interval)
	}
}

// reap removes all cache entries that have expired based on the provided interval.
// An entry is considered expired if its creation time is before (now - interval).
//
// Parameters:
//   - interval: The time duration used to determine if an entry has expired
func (c *Cache) reap(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	timeAgo := time.Now().UTC().Add(-interval)
	for k, v := range c.cache {
		if v.createdAt.Before(timeAgo) {
			delete(c.cache, k)
		}
	}
}
