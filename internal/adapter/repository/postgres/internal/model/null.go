package model

import (
	"database/sql"
)

func NullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NullIntPositive(i int32) sql.NullInt32 {
	if i <= 0 {
		return sql.NullInt32{}
	}

	return sql.NullInt32{
		Int32: i,
		Valid: true,
	}
}
