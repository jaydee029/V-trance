package utils

import "github.com/jackc/pgx/v5/pgtype"

func ToText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  s != "",
	}
}
