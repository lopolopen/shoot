package mapper

import (
	"database/sql"
	"time"
)

type SQLMapper struct{}

func (m *SQLMapper) TimePtrToNullTime(t *time.Time) sql.Null[time.Time] {
	if t == nil {
		return sql.Null[time.Time]{}
	}
	return sql.Null[time.Time]{
		V:     *t,
		Valid: true,
	}
}
