package bootstrap

import (
	"context"
	"log"

	"github.com/iamKienb/shopify-go-platform/cache/redis"
	"github.com/iamKienb/shopify-go-platform/config"
	"github.com/iamKienb/shopify-go-platform/database/postgres"
	"github.com/iamKienb/shopify-go-platform/messaging/kafka"
)

type Bootstrap struct {
	cfg   *config.Config
	pg    *postgres.Client
	redis *redis.Client
	kafka *kafka.Client
}

func NewBootstrap() *Bootstrap {
	cfg := config.Load()

	return &Bootstrap{
		cfg: cfg,
	}
}

func (a *Bootstrap) Start(ctx context.Context) {
	var err error

	a.pg, err = postgres.New(a.cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	a.redis, err = redis.New(a.cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}

	// 	a.kafka, err = kafka.New(a.cfg.Kafka)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
}

func (a *Bootstrap) Stop() {
	if a.pg != nil {
		a.pg.Close()
	}
	if a.redis != nil {
		a.redis.Close()
	}
	// if a.kafka != nil {
	// 	a.kafka.Close()
	// }
}
