package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	t.Run("enabled cache", func(t *testing.T) {
		c, err := NewCache(true)
		if err != nil {
			t.Fatalf("NewCache failed: %v", err)
		}

		if !c.IsEnabled() {
			t.Error("Cache should be enabled")
		}

		if c.GetCacheDir() == "" {
			t.Error("Cache directory should not be empty")
		}

		// Check if cache directory exists
		if _, err := os.Stat(c.GetCacheDir()); os.IsNotExist(err) {
			t.Error("Cache directory should exist")
		}
	})

	t.Run("disabled cache", func(t *testing.T) {
		c, err := NewCache(false)
		if err != nil {
			t.Fatalf("NewCache failed: %v", err)
		}

		if c.IsEnabled() {
			t.Error("Cache should be disabled")
		}
	})
}

func TestCacheSetAndGet(t *testing.T) {
	c, err := NewCache(true)
	if err != nil {
		t.Fatalf("NewCache failed: %v", err)
	}

	testURL := "https://example.com/test.lua"
	testContent := "-- Test content\nprint('Hello, World!')"

	// Set cache
	if err := c.Set(testURL, testContent); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get cache
	content, found, err := c.Get(testURL)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if !found {
		t.Error("Cache should be found")
	}

	if content != testContent {
		t.Errorf("Expected content %q, got %q", testContent, content)
	}

	// Clean up
	c.Clear()
}

func TestCacheExpiry(t *testing.T) {
	c, err := NewCache(true)
	if err != nil {
		t.Fatalf("NewCache failed: %v", err)
	}

	testURL := "https://example.com/expiry-test.lua"
	testContent := "-- Expiry test"

	// Set cache
	if err := c.Set(testURL, testContent); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Manually modify file time to simulate expiry
	cacheKey := c.generateCacheKey(testURL)
	cachePath := filepath.Join(c.GetCacheDir(), cacheKey)

	oldTime := time.Now().Add(-25 * time.Hour) // More than 24 hours ago
	if err := os.Chtimes(cachePath, oldTime, oldTime); err != nil {
		t.Fatalf("Failed to modify file time: %v", err)
	}

	// Try to get expired cache
	_, found, err := c.Get(testURL)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if found {
		t.Error("Expired cache should not be found")
	}

	// Clean up
	c.Clear()
}

func TestCacheDisabled(t *testing.T) {
	c, err := NewCache(false)
	if err != nil {
		t.Fatalf("NewCache failed: %v", err)
	}

	testURL := "https://example.com/disabled-test.lua"
	testContent := "-- Disabled test"

	// Set should not error even when disabled
	if err := c.Set(testURL, testContent); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get should return not found
	_, found, err := c.Get(testURL)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if found {
		t.Error("Cache should not be found when disabled")
	}
}

func TestCacheClear(t *testing.T) {
	c, err := NewCache(true)
	if err != nil {
		t.Fatalf("NewCache failed: %v", err)
	}

	// Add multiple items to cache
	urls := []string{
		"https://example.com/test1.lua",
		"https://example.com/test2.lua",
		"https://example.com/test3.lua",
	}

	for _, url := range urls {
		if err := c.Set(url, "test content"); err != nil {
			t.Fatalf("Set failed: %v", err)
		}
	}

	// Clear cache
	if err := c.Clear(); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// Verify all items are gone
	for _, url := range urls {
		_, found, err := c.Get(url)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if found {
			t.Errorf("Cache for %s should not be found after clear", url)
		}
	}
}

func TestGenerateCacheKey(t *testing.T) {
	c, _ := NewCache(true)

	url1 := "https://example.com/test.lua"
	url2 := "https://example.com/test.lua"
	url3 := "https://example.com/other.lua"

	key1 := c.generateCacheKey(url1)
	key2 := c.generateCacheKey(url2)
	key3 := c.generateCacheKey(url3)

	// Same URLs should generate same keys
	if key1 != key2 {
		t.Error("Same URLs should generate same cache keys")
	}

	// Different URLs should generate different keys
	if key1 == key3 {
		t.Error("Different URLs should generate different cache keys")
	}

	// Keys should end with .lua
	if filepath.Ext(key1) != ".lua" {
		t.Error("Cache key should have .lua extension")
	}

	// Clean up
	c.Clear()
}
