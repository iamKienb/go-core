package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	rdb *redis.Client
}

func NewRedisService(rdb *redis.Client) RedisAbstract {
	return &RedisService{rdb: rdb}
}

func (r *RedisService) Get(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

func (r *RedisService) Set(ctx context.Context, key string, value any, ttlSeconds int) error {
	return r.rdb.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *RedisService) Del(ctx context.Context, keys ...string) error {
	return r.rdb.Del(ctx, keys...).Err()

}

func (r *RedisService) Exists(ctx context.Context, key string) (bool, error) {
	res, err := r.rdb.Exists(ctx, key).Result()
	return res > 0, err

}

func (r *RedisService) Expire(ctx context.Context, key string, ttlSeconds int) error {
	return r.rdb.Expire(
		ctx,
		key,
		time.Duration(ttlSeconds)*time.Second,
	).Err()
}
