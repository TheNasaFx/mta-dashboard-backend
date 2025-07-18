package model

type Market struct {
	ID                  uint    `json:"id"`
	OpTypeName          string  `json:"op_type_name"`
	DistCode            string  `json:"dist_code"`
	KhoCode             string  `json:"kho_code"`
	StorName            string  `json:"stor_name"`
	StorFloor           string  `json:"stor_floor"`
	MrchRegno           string  `json:"mrch_regno"`
	PayCenterPropertyID uint    `json:"pay_center_property_id"`
	PayCenterID         uint    `json:"pay_center_id"`
	Lat                 float64 `json:"lat"`
	Lng                 float64 `json:"lng"`
	BuildFloor          *int    `json:"build_floor"`
}
