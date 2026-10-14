package redisx

import (
	"context"
	"fmt"
	"time"

	configx "github.com/iamKienb/go-core/config"
	redisotel "github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type RedisX struct {
	client *redis.Client
}

func New(ctx context.Context, cfg configx.RedisConfig) (RedisXService, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.Db,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: 10,
	})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		rdb.Close()
		return nil, fmt.Errorf("redis tracing instrument failed: %w", err)
	}
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		rdb.Close()
		return nil, fmt.Errorf("redis metrics instrument failed: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		rdb.Close()
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &RedisX{client: rdb}, nil
}
