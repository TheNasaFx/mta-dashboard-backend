package repository

import (
	"dashboard-backend/database"
	"dashboard-backend/database/model"
	"fmt"
)

// GetEbarimtByPin returns ebarimt info by PIN
func GetEbarimtByPin(pin string) (*model.PayMarketEbarimt, error) {
	if database.DB == nil {
		database.MustConnect()
	}
	query := `SELECT ID, PIN, COUNT_RECEIPT FROM GPS.PAY_MARKET_BARIMT WHERE TRIM(UPPER(PIN)) = TRIM(UPPER(:1))`
	row := database.DB.QueryRow(query, pin)
	var e model.PayMarketEbarimt
	err := row.Scan(&e.ID, &e.Pin, &e.CountReceipt)
	if err != nil {
		return nil, fmt.Errorf("DB query error: %w", err)
	}
	return &e, nil
}
