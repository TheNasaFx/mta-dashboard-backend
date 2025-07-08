package repository

import (
	"dashboard-backend/database"
	"dashboard-backend/database/model"
	"fmt"
)

// GetEbarimtByPin returns aggregated ebarimt info by PIN
func GetEbarimtByPin(pin string) (*model.PayMarketEbarimt, error) {
	if database.DB == nil {
		database.MustConnect()
	}

	// Use SUM to aggregate all COUNT_RECEIPT values for the PIN
	query := `SELECT 
		MIN(ID) as ID, 
		:1 as PIN, 
		SUM(COUNT_RECEIPT) as TOTAL_COUNT_RECEIPT 
	FROM GPS.PAY_MARKET_BARIMT 
	WHERE TRIM(UPPER(PIN)) = TRIM(UPPER(:1))`

	row := database.DB.QueryRow(query, pin)
	var e model.PayMarketEbarimt
	err := row.Scan(&e.ID, &e.Pin, &e.CountReceipt)
	if err != nil {
		return nil, fmt.Errorf("DB query error: %w", err)
	}
	return &e, nil
}
