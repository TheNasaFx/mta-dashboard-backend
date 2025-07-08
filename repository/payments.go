package repository

import (
	"dashboard-backend/database"
	"dashboard-backend/database/model"
	"fmt"
)

// GetPaymentsByPin returns payment information by PIN
func GetPaymentsByPin(pin string) (*model.PaymentsSummary, error) {
	if database.DB == nil {
		database.MustConnect()
	}

	// Get summary data (earliest and latest dates, total amount)
	summaryQuery := `SELECT 
		MIN(PAID_DATE) as EARLIEST_DATE,
		MAX(PAID_DATE) as LATEST_DATE,
		SUM(AMOUNT) as TOTAL_AMOUNT,
		COUNT(*) as PAYMENT_COUNT
	FROM GPS.V_E_TUB_PAYMENTS 
	WHERE TRIM(UPPER(PIN)) = TRIM(UPPER(:1))`

	row := database.DB.QueryRow(summaryQuery, pin)
	var summary model.PaymentsSummary
	err := row.Scan(&summary.EarliestDate, &summary.LatestDate, &summary.TotalAmount, &summary.PaymentCount)
	if err != nil {
		return nil, fmt.Errorf("summary query error: %w", err)
	}

	// Get detailed payment records
	detailQuery := `SELECT 
		TAX_TYPE_NAME,
		INVOICE_NO,
		BRANCH_NAME,
		PAID_DATE,
		AMOUNT,
		DESCRIPTION,
		ENTITY_NAME,
		PAYMENT_METHOD_NAME,
		TAX_TYPE_CODE,
		BRANCH_CODE
	FROM GPS.V_E_TUB_PAYMENTS 
	WHERE TRIM(UPPER(PIN)) = TRIM(UPPER(:1))
	ORDER BY PAID_DATE DESC`

	rows, err := database.DB.Query(detailQuery, pin)
	if err != nil {
		return nil, fmt.Errorf("detail query error: %w", err)
	}
	defer rows.Close()

	var payments []model.PaymentDetail
	for rows.Next() {
		var payment model.PaymentDetail
		err := rows.Scan(
			&payment.TaxTypeName,
			&payment.InvoiceNo,
			&payment.BranchName,
			&payment.PaidDate,
			&payment.Amount,
			&payment.Description,
			&payment.EntityName,
			&payment.PaymentMethodName,
			&payment.TaxTypeCode,
			&payment.BranchCode,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		payments = append(payments, payment)
	}

	summary.Payments = payments
	summary.PIN = pin
	return &summary, nil
}
