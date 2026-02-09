package postgres

import (
	"context"
)

func PgWithTx(
	ctx context.Context,
	pool PgPool,
	fn func(ctx context.Context, tx PgTx) error,
) error {
	tx, err := pool.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
