package model

import "database/sql"

type PayCenterLocation struct {
	PayCenterID sql.NullInt64   `db:"PAY_CENTER_ID"`
	LNG         sql.NullFloat64 `db:"LNG"`
	LAT         sql.NullFloat64 `db:"LAT"`
}
