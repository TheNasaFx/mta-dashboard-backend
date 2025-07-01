package repository

import (
	"dashboard-backend/database"
	"dashboard-backend/database/model"
	"fmt"
)

func GetMarketsByOrgID(orgID uint) ([]model.Market, error) {
	if database.DB == nil {
		database.MustConnect()
	}
	query := `SELECT ID, OP_TYPE_NAME, DIST_CODE, KHO_CODE, STOR_NAME, STOR_FLOOR, MRCH_REGNO, PAY_CENTER_PROPERTY_ID, PAY_CENTER_ID, LAT, LNG
			  FROM GPS.PAY_MARKET
			  WHERE PAY_CENTER_ID = :orgID`
	rows, err := database.DB.Query(query, orgID)
	if err != nil {
		return nil, fmt.Errorf("DB query error: %w", err)
	}
	defer rows.Close()

	var markets []model.Market
	for rows.Next() {
		var m model.Market
		err := rows.Scan(
			&m.ID,
			&m.OpTypeName,
			&m.DistCode,
			&m.KhoCode,
			&m.StorName,
			&m.StorFloor,
			&m.MrchRegno,
			&m.PayCenterPropertyID,
			&m.PayCenterID,
			&m.Lat,
			&m.Lng,
		)
		if err != nil {
			return nil, fmt.Errorf("Scan error: %w", err)
		}
		markets = append(markets, m)
	}
	return markets, nil
}
