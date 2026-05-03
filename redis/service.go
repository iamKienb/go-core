package redisx

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

func (x *RedisX) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return x.client.Set(ctx, key, data, ttl).Err()
}

func (x *RedisX) Get(ctx context.Context, key string, dest any) error {
	val, err := x.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return redis.Nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (x *RedisX) Exists(ctx context.Context, key string) (bool, error) {
	n, err := x.client.Exists(ctx, key).Result()
	return n > 0, err
}

func (x *RedisX) Delete(ctx context.Context, key string) error {
	return x.client.Del(ctx, key).Err()
}

func (x *RedisX) Incr(ctx context.Context, key string) (int64, error) {
	return x.client.Incr(ctx, key).Result()
}

func (x *RedisX) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return x.client.Expire(ctx, key, ttl).Err()
}

func (x *RedisX) Lock(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	return x.client.SetNX(ctx, key, value, ttl).Result()
}

func (x *RedisX) Unlock(ctx context.Context, key string, value string) error {
	var script = redis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
        else
            return 0
        end
	`)
	res, err := script.Run(ctx, x.client, []string{key}, value).Int()
	if err != nil {
		return err
	}
	if res == 0 {
		return errors.New("could not unlock: key not found or value mismatch")
	}
	return nil
}

func (x *RedisX) GetClient() *redis.Client {
	return x.client
}

func (x *RedisX) Close() error {
	return x.client.Close()
}
