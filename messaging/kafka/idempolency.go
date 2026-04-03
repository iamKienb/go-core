package kafka

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Idempotency struct {
	client *redis.Client
	ttl    time.Duration
}

func NewIdempotency(rdb *redis.Client, ttl time.Duration) *Idempotency {
	return &Idempotency{
		client: rdb,
		ttl:    ttl,
	}
}

func (i *Idempotency) CheckAndSet(ctx context.Context, key string) bool {
	ok, _ := i.client.SetNX(ctx, key, "1", i.ttl).Result()
	return ok
}
