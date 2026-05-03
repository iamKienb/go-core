package redisx

import (
	"context"
	"fmt"
	"time"

	configx "github.com/iamKienb/shopify-go-platform/config"
	"github.com/redis/go-redis/v9"
)

type RedisX struct {
	client *redis.Client
}

func New(cfg configx.RedisConfig) (RedisXService, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.Db,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: 10,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &RedisX{client: rdb}, nil
}
