package model

import "database/sql"

type Property struct {
	ID           uint            `json:"id"`
	PayCenterID  uint            `json:"pay_center_id"`
	UpdatedDate  sql.NullString  `json:"updated_date"`
	PropertyType sql.NullString  `json:"property_type"`
	OwnerRegno   sql.NullString  `json:"owner_regno"`
	PropertySize sql.NullFloat64 `json:"property_size"`
	RentAmount   sql.NullFloat64 `json:"rent_amount"`
}
