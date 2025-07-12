// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig holds configuration for Redis cache backend
type RedisConfig struct {
	Addr       string
	Password   string
	DB         int
	KeyPrefix  string
	DefaultTTL time.Duration
}

// RedisCacheBackend implements CacheBackend using Redis
type RedisCacheBackend[K comparable, V any] struct {
	client     *redis.Client
	ctx        context.Context
	keyPrefix  string
	defaultTTL time.Duration
}

// NewRedisCacheBackend creates a new Redis cache backend
func NewRedisCacheBackend[K comparable, V any](cfg RedisConfig) (*RedisCacheBackend[K, V], error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &RedisCacheBackend[K, V]{
		client:     client,
		ctx:        context.Background(),
		keyPrefix:  cfg.KeyPrefix,
		defaultTTL: cfg.DefaultTTL,
	}, nil
}

func (r *RedisCacheBackend[K, V]) redisKey(key K) string {
	keyStr := toString(key)
	if r.keyPrefix != "" {
		return r.keyPrefix + keyStr
	}
	return keyStr
}

func (r *RedisCacheBackend[K, V]) Get(key K) (V, bool) {
	var zero V
	data, err := r.client.Get(r.ctx, r.redisKey(key)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return zero, false
		}
		return zero, false
	}

	var value V
	if err := json.Unmarshal(data, &value); err != nil {
		return zero, false
	}
	return value, true
}

func (r *RedisCacheBackend[K, V]) Set(key K, value V) {
	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	r.client.Set(r.ctx, r.redisKey(key), data, r.defaultTTL)
}

func (r *RedisCacheBackend[K, V]) Remove(key K) {
	r.client.Del(r.ctx, r.redisKey(key))
}

func (r *RedisCacheBackend[K, V]) Purge() {
	// WARNING: This deletes all keys in the current DB. Use with caution.
	r.client.FlushDB(r.ctx)
}

func (r *RedisCacheBackend[K, V]) Len() int {
	size, err := r.client.DBSize(r.ctx).Result()
	if err != nil {
		return 0
	}
	return int(size)
}

// toString converts any comparable type to string for Redis key
func toString[K comparable](key K) string {
	// For now, we'll use JSON marshaling as a generic approach
	// This could be optimized for specific types if needed
	data, _ := json.Marshal(key)
	return string(data)
}
