package conv

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func UUID(id [16]byte) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func Text(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}

	return pgtype.Text{String: *s, Valid: true}
}

func Number(n *int64) pgtype.Int8 {
	if n == nil {
		return pgtype.Int8{Valid: false}
	}

	return pgtype.Int8{Int64: *n, Valid: true}
}

func Date(d *time.Time) pgtype.Date {
	if d == nil {
		return pgtype.Date{Valid: false}
	}

	return pgtype.Date{Time: *d, Valid: true}
}

func TimeStampZ(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}

	return pgtype.Timestamptz{Time: *t, Valid: true}
}
