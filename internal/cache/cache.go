package cache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	cacheDirName = ".lua-bundler-cache"
	cacheExpiry  = 24 * time.Hour // Cache expires after 24 hours
)

type Cache struct {
	cacheDir string
	enabled  bool
}

// NewCache creates a new cache instance
func NewCache(enabled bool) (*Cache, error) {
	if !enabled {
		return &Cache{enabled: false}, nil
	}

	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	cacheDir := filepath.Join(homeDir, cacheDirName)

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &Cache{
		cacheDir: cacheDir,
		enabled:  true,
	}, nil
}

// generateCacheKey creates a unique cache key from URL
func (c *Cache) generateCacheKey(url string) string {
	hash := md5.Sum([]byte(url))
	return hex.EncodeToString(hash[:]) + ".lua"
}

// Get retrieves content from cache if it exists and is not expired
func (c *Cache) Get(url string) (string, bool, error) {
	if !c.enabled {
		return "", false, nil
	}

	cacheKey := c.generateCacheKey(url)
	cachePath := filepath.Join(c.cacheDir, cacheKey)

	// Check if cache file exists
	info, err := os.Stat(cachePath)
	if os.IsNotExist(err) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	// Check if cache is expired
	if time.Since(info.ModTime()) > cacheExpiry {
		// Delete expired cache
		os.Remove(cachePath)
		return "", false, nil
	}

	// Read cache file
	content, err := os.ReadFile(cachePath)
	if err != nil {
		return "", false, err
	}

	return string(content), true, nil
}

// Set stores content in cache
func (c *Cache) Set(url string, content string) error {
	if !c.enabled {
		return nil
	}

	cacheKey := c.generateCacheKey(url)
	cachePath := filepath.Join(c.cacheDir, cacheKey)

	// Write content to cache file
	if err := os.WriteFile(cachePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write cache: %w", err)
	}

	return nil
}

// Clear removes all cached files
func (c *Cache) Clear() error {
	if !c.enabled {
		return nil
	}

	// Remove all files in cache directory
	entries, err := os.ReadDir(c.cacheDir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			cachePath := filepath.Join(c.cacheDir, entry.Name())
			if err := os.Remove(cachePath); err != nil {
				return fmt.Errorf("failed to remove cache file %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

// GetCacheDir returns the cache directory path
func (c *Cache) GetCacheDir() string {
	return c.cacheDir
}

// IsEnabled returns whether cache is enabled
func (c *Cache) IsEnabled() bool {
	return c.enabled
}
