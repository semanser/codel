package database

import (
	"database/sql"
)

func StringToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}
