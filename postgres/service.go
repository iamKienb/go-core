package pgx

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

func (x *PGX) GetPool() *pgxpool.Pool {
	return x.pool
}

func (x *PGX) Close() {
	x.pool.Close()
}
