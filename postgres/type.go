package pgx

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXService interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
	GetPool() *pgxpool.Pool
	Close()
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
