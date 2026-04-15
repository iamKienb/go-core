package redisx

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/iamKienb/shopify-go-platform/config"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	Conn *redis.Client
}

func New(cfg config.RedisConfig) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.Database,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: 10,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	log.Printf("Redis connected successfully to port: %d", cfg.Port)

	return &Client{Conn: rdb}, nil
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
