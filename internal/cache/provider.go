package cache

import (
	"context"
	"time"
)

type Provider interface {
	Get(ctx context.Context, key string) ([]byte, bool)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration)
	Invalidate(ctx context.Context, key string)
}
