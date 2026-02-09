package redis

import (
	"context"
)

type Locker interface {
	Lock(ctx context.Context, key string, ttlSeconds int) (bool, error)
	Unlock(ctx context.Context, key string) error
}
