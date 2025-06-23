package model

const TableNameOrg = "pay_center.org"

type Org struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	// Code      string        `json:"code"`
	// Type      string        `json:"type"`
	// Address   string        `json:"address"`
	// Phone     string        `json:"phone"`
	// Email     string        `json:"email"`
	// Status    string        `json:"status"`
	// ParentID  sql.NullInt64 `json:"parent_id"`
	// CreatedAt time.Time     `json:"created_at"`
	// UpdatedAt time.Time     `json:"updated_at"`
	// DeletedAt sql.NullTime  `json:"deleted_at"`
}

func (*Org) TableName() string {
	return TableNameOrg
}
