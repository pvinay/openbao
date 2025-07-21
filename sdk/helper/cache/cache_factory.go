package cache

// NewCacheBackend creates a cache backend based on environment variables
func NewCacheBackend[K comparable, V any](size int) (CacheBackend[K, V], error) {
	cacheType := GetCacheBackendType()

	switch cacheType {
	case CacheBackendRedis:
		return NewRedisCacheBackend[K, V](RedisConfig{
			Addr:       GetEnvOrDefault("REDIS_ADDR", "localhost:6379"),
			Password:   GetEnvOrDefault("REDIS_PASSWORD", ""),
			DB:         GetEnvIntOrDefault("REDIS_DB", 0),
			KeyPrefix:  GetEnvOrDefault("REDIS_KEY_PREFIX", ""),
			DefaultTTL: 0, // No TTL by default
		})
	case CacheBackendLRU:
		fallthrough
	default:
		return NewLRUCacheBackend[K, V](size)
	}
}
