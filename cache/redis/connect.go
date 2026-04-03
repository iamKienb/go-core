package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/iamKienb/shopify-go-platform/config"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	redis *redis.Client
}

func New(cfg config.RedisConfig) (*Client, error) {
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: cfg.Password,
		Username: cfg.Username,
		DB:       cfg.Database,
		PoolSize: cfg.PoolSize,
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis cluster ping failed:", err)
	}
	fmt.Println("Redis cluster connected:", pong)

	return &Client{
		redis: rdb,
	}, nil

}

func (c *Client) Close() {
	c.redis.Close()
}
