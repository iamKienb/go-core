package postgres

import (
	"context"
	"fmt"

	"github.com/iamKienb/shopify-go-platform/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	pool *pgxpool.Pool
}

func New(cfg config.PostgresConfig) (*Client, error) {

	dns := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Db,
	)

	poolConfig, err := pgxpool.ParseConfig(dns)

	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = cfg.MaxOpenConns
	poolConfig.MinConns = cfg.MaxIdleConns
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)

	if err != nil {
		return nil, err
	}

	fmt.Println("Set Pool Pg Success")

	return &Client{pool: pool}, nil

}

func (c *Client) Close() {
	c.pool.Close()
}

func (c *Client) GetPool() *pgxpool.Pool {
	return c.pool
}
