package domain

import (
	"database/sql"
	"time"
)

// Helper functions to create sql.Null* types
func SQLNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func SQLNullInt32(i int) sql.NullInt32 {
	return sql.NullInt32{Int32: int32(i), Valid: true}
}

func SQLNullTime(t time.Time) sql.NullTime {
    if t.IsZero() {
        return sql.NullTime{Valid: false}
    }
    return sql.NullTime{Time: t, Valid: true}
}
