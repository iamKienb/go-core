package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {

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

	fmt.Println("Set Pool Pg Success")

	poolConfig.MaxConns = cfg.MaxOpenConns
	poolConfig.MinConns = cfg.MaxIdleConns
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime

	return pgxpool.NewWithConfig(ctx, poolConfig)

}
