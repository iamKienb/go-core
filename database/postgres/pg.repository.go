package postgres

import "context"

type PgRepository struct {
	p PgPool
}

func NewPgRepository(p PgPool) PgDB {
	return &PgRepository{p: p}
}

func (a *PgRepository) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := a.p.Exec(ctx, sql, args...)
	return err

}

func (a *PgRepository) Query(ctx context.Context, sql string, args ...any) (Rows, error) {
	return a.p.Query(ctx, sql, args...)

}

func (a *PgRepository) QueryRow(ctx context.Context, sql string, args ...any) Row {
	return a.p.QueryRow(ctx, sql, args...)
}
