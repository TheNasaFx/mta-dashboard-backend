package repository

import (
	"dashboard-backend/database"
	"dashboard-backend/database/model"
	"database/sql"
	"fmt"
)

func GetMarketsByOrgID(orgID uint) ([]model.Market, error) {
	if database.DB == nil {
		database.MustConnect()
	}
	query := `SELECT pm.ID, pm.OP_TYPE_NAME, pm.DIST_CODE, pm.KHO_CODE, pm.STOR_NAME, pm.STOR_FLOOR, pm.MRCH_REGNO, pm.PAY_CENTER_PROPERTY_ID, pm.PAY_CENTER_ID, pm.LAT, pm.LNG, pc.BUILD_FLOOR
			  FROM GPS.PAY_MARKET pm
			  LEFT JOIN GPS.PAY_CENTER pc ON pm.PAY_CENTER_ID = pc.ID
			  WHERE pm.PAY_CENTER_ID = :orgID`
	rows, err := database.DB.Query(query, orgID)
	if err != nil {
		return nil, fmt.Errorf("DB query error: %w", err)
	}
	defer rows.Close()

	var markets []model.Market
	for rows.Next() {
		var m model.Market
		var buildFloor sql.NullInt64
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
			&buildFloor,
		)
		if err != nil {
			return nil, fmt.Errorf("Scan error: %w", err)
		}
		if buildFloor.Valid {
			buildFloorInt := int(buildFloor.Int64)
			m.BuildFloor = &buildFloorInt
		}
		markets = append(markets, m)
	}
	return markets, nil
}
