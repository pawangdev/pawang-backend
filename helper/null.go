package helper

import "database/sql"

func IntToSQLNullInt32(id int32) sql.NullInt32 {
	if id != 0 {
		return sql.NullInt32{
			Int32: id,
			Valid: true,
		}
	}

	return sql.NullInt32{}
}
