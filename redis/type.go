package redisx

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string, dest any) error
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error

	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
}

type Locker interface {
	Lock(ctx context.Context, key string, value string, ttl time.Duration) (bool, error)
	Unlock(ctx context.Context, key string, value string) error
}

type RedisXService interface {
	Cache
	Locker
	GetClient() *redis.Client
	Close() error
}
