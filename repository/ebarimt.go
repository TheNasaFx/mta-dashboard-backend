package repository

import (
	"dashboard-backend/database"
	"dashboard-backend/database/model"
	"fmt"
)

// GetEbarimtByPin returns aggregated ebarimt info by PIN from the new table structure
func GetEbarimtByPin(pin string) (*model.PayMarketEbarimt, error) {
	if database.DB == nil {
		database.MustConnect()
	}

	// Query the new V_E_TUB_PAY_MARKET_EBARIMT table with new columns
	query := `SELECT 
		MIN(ID) as ID, 
		:1 as PIN, 
		SUM(CNT_3) as CNT_3,
		SUM(CNT_30) as CNT_30
	FROM GPS.V_E_TUB_PAY_MARKET_EBARIMT 
	WHERE TRIM(UPPER(MRCH_REGNO)) = TRIM(UPPER(:1))`

	row := database.DB.QueryRow(query, pin)
	var e model.PayMarketEbarimt
	err := row.Scan(&e.ID, &e.Pin, &e.Cnt3, &e.Cnt30)
	if err != nil {
		return nil, fmt.Errorf("DB query error: %w", err)
	}

	// Set CountReceipt for backward compatibility (use CNT_3 as default)
	e.CountReceipt = e.Cnt3

	return &e, nil
}
