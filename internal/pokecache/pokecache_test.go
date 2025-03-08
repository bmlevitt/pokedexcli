// This file contains tests for the pokecache package.
// It verifies the functionality of the cache implementation, including
// creation, adding/retrieving values, and automatic expiration of entries.
package pokecache

import (
	"testing"
	"time"
)

// TestCreateCache verifies that a new cache can be created successfully.
// It checks that the internal map is properly initialized and not nil.
func TestCreateCache(t *testing.T) {
	cache := NewCache(time.Millisecond)
	if cache.cache == nil {
		t.Error("cache is nil")
	}
}

// TestAddGetCache verifies that values can be added to the cache and retrieved correctly.
// It tests multiple key-value pairs, including an empty key, to ensure the cache
// correctly stores and retrieves all values.
func TestAddGetCache(t *testing.T) {
	cache := NewCache(time.Millisecond)

	cases := []struct {
		inputKey string
		inputVal []byte
	}{
		{
			inputKey: "key1",
			inputVal: []byte("val1"),
		},
		{
			inputKey: "key2",
			inputVal: []byte("val2"),
		},
		{
			inputKey: "",
			inputVal: []byte("val3"),
		},
	}

	for _, cas := range cases {
		cache.Add(cas.inputKey, cas.inputVal)
		actual, ok := cache.Get(cas.inputKey)
		if !ok {
			t.Errorf("input key %v not found in cache", cas.inputKey)
		}
		if string(actual) != string(cas.inputVal) {
			t.Errorf("input value %v does not match actual value %v", cas.inputVal, string(actual))
		}
	}
}

// TestReap verifies that the cache's automatic cleanup mechanism works correctly.
// It adds an entry to the cache, waits for longer than the expiration interval,
// and then checks that the entry has been automatically removed ("reaped").
func TestReap(t *testing.T) {
	interval := time.Millisecond * 10
	cache := NewCache(interval)
	keyOne := "key1"
	cache.Add(keyOne, []byte("val1"))
	time.Sleep(interval * 2)

	_, ok := cache.Get(keyOne)
	if ok {
		t.Errorf("%s should have been reaped", keyOne)
	}
}
