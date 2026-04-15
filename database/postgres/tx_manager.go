package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

type pgxTxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) TxManager {
	return &pgxTxManager{
		pool: pool,
	}
}

func (m *pgxTxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.pool.Begin(ctx)

	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		_ = tx.Rollback(txCtx)
		return err
	}

	if err := tx.Commit(txCtx); err != nil {
		return fmt.Errorf("Commit transaction: %w", err)
	}
	return nil

}

func ExtractTx(ctx context.Context) pgx.Tx {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	if !ok {
		return nil
	}
	return tx
}
