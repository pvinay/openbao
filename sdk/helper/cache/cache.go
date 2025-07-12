// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cache

import (
	"os"
	"strconv"
	"strings"
)

// CacheBackend defines a generic interface for cache implementations
type CacheBackend[K comparable, V any] interface {
	// Get retrieves a value from the cache by key
	Get(key K) (V, bool)

	// Set stores a value in the cache
	Set(key K, value V)

	// Remove removes a value from the cache by key
	Remove(key K)

	// Purge clears all entries from the cache
	Purge()

	// Len returns the number of entries in the cache
	Len() int
}

// CacheBackendType represents the type of cache backend
type CacheBackendType string

const (
	CacheBackendLRU   CacheBackendType = "lru"
	CacheBackendRedis CacheBackendType = "redis"
)

// Environment variable to control cache backend
const CacheBackendEnvVar = "CACHE_BACKEND"

// GetCacheBackendType returns the cache backend type from environment variable
func GetCacheBackendType() CacheBackendType {
	cacheType := strings.ToLower(strings.TrimSpace(os.Getenv(CacheBackendEnvVar)))
	switch cacheType {
	case string(CacheBackendRedis):
		return CacheBackendRedis
	case string(CacheBackendLRU), "":
		return CacheBackendLRU
	default:
		return CacheBackendLRU
	}
}

// Helper functions for environment variables
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
