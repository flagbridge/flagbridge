package cache

import (
	"context"
	"log/slog"
	"time"

	"github.com/dgraph-io/ristretto"
)

type MemoryCache struct {
	cache *ristretto.Cache
}

func NewMemoryCache() (*MemoryCache, error) {
	c, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e5,
		MaxCost:     1 << 24, // 16MB
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}
	slog.Info("initialized in-memory cache")
	return &MemoryCache{cache: c}, nil
}

func (m *MemoryCache) Get(_ context.Context, key string) ([]byte, bool) {
	val, found := m.cache.Get(key)
	if !found {
		return nil, false
	}
	b, ok := val.([]byte)
	return b, ok
}

func (m *MemoryCache) Set(_ context.Context, key string, value []byte, ttl time.Duration) {
	m.cache.SetWithTTL(key, value, int64(len(value)), ttl)
}

func (m *MemoryCache) Invalidate(_ context.Context, key string) {
	m.cache.Del(key)
}
