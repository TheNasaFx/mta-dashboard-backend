package model

const TableNameOrg = "pay_center.org"

type Org struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	OfficeCode string  `json:"office_code"`
	Regno      string  `json:"regno"`
	KhoCode    string  `json:"kho_code"`
	BuildFloor string  `json:"build_floor"`
	Address    string  `json:"address"`
	Lng        float64 `json:"lng"`
	Lat        float64 `json:"lat"`
}

func (*Org) TableName() string {
	return TableNameOrg
}
