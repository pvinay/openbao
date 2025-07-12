// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cache

import (
	lru "github.com/hashicorp/golang-lru/v2"
)

// LRUCacheBackend implements CacheBackend using the LRU cache
type LRUCacheBackend[K comparable, V any] struct {
	cache *lru.TwoQueueCache[K, V]
}

// NewLRUCacheBackend creates a new LRU cache backend
func NewLRUCacheBackend[K comparable, V any](size int) (*LRUCacheBackend[K, V], error) {
	cache, err := lru.New2Q[K, V](size)
	if err != nil {
		return nil, err
	}
	return &LRUCacheBackend[K, V]{cache: cache}, nil
}

func (l *LRUCacheBackend[K, V]) Get(key K) (V, bool) {
	return l.cache.Get(key)
}

func (l *LRUCacheBackend[K, V]) Set(key K, value V) {
	l.cache.Add(key, value)
}

func (l *LRUCacheBackend[K, V]) Remove(key K) {
	l.cache.Remove(key)
}

func (l *LRUCacheBackend[K, V]) Purge() {
	l.cache.Purge()
}

func (l *LRUCacheBackend[K, V]) Len() int {
	return l.cache.Len()
}
