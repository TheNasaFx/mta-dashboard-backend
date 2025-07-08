package model

import (
	"database/sql"
)

// PaymentDetail represents individual payment record
type PaymentDetail struct {
	TaxTypeName       sql.NullString  `json:"tax_type_name"`
	InvoiceNo         sql.NullString  `json:"invoice_no"`
	BranchName        sql.NullString  `json:"branch_name"`
	PaidDate          sql.NullTime    `json:"paid_date"`
	Amount            sql.NullFloat64 `json:"amount"`
	Description       sql.NullString  `json:"description"`
	EntityName        sql.NullString  `json:"entity_name"`
	PaymentMethodName sql.NullString  `json:"payment_method_name"`
	TaxTypeCode       sql.NullString  `json:"tax_type_code"`
	BranchCode        sql.NullString  `json:"branch_code"`
}

// PaymentsSummary represents aggregated payment information
type PaymentsSummary struct {
	PIN          string          `json:"pin"`
	EarliestDate sql.NullTime    `json:"earliest_date"`
	LatestDate   sql.NullTime    `json:"latest_date"`
	TotalAmount  sql.NullFloat64 `json:"total_amount"`
	PaymentCount int             `json:"payment_count"`
	Payments     []PaymentDetail `json:"payments"`
}
