package repository

import (
	"dashboard-backend/database/model"
	"database/sql"
)

func GetPayCenterLocationsByPayCenterID(db *sql.DB, payCenterID int64) ([]model.PayCenterLocation, error) {
	rows, err := db.Query("SELECT PAY_CENTER_ID, LNG, LAT FROM GPS.PAY_CENTER_LOCATION WHERE PAY_CENTER_ID = :1", payCenterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.PayCenterLocation
	for rows.Next() {
		var loc model.PayCenterLocation
		if err := rows.Scan(&loc.PayCenterID, &loc.LNG, &loc.LAT); err != nil {
			return nil, err
		}
		results = append(results, loc)
	}
	if results == nil {
		results = make([]model.PayCenterLocation, 0)
	}
	return results, nil
}

// New: get pay_center_location by regno (join pay_center)
func GetPayCenterLocationsByRegno(db *sql.DB, regno string) ([]model.PayCenterLocation, error) {
	query := `SELECT pcl.PAY_CENTER_ID, pcl.LNG, pcl.LAT
			  FROM GPS.PAY_CENTER_LOCATION pcl
			  JOIN GPS.PAY_CENTER o ON pcl.PAY_CENTER_ID = o.ID
			  WHERE o.REGNO = :1`
	rows, err := db.Query(query, regno)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []model.PayCenterLocation
	for rows.Next() {
		var loc model.PayCenterLocation
		if err := rows.Scan(&loc.PayCenterID, &loc.LNG, &loc.LAT); err != nil {
			return nil, err
		}
		results = append(results, loc)
	}
	if results == nil {
		results = make([]model.PayCenterLocation, 0)
	}
	return results, nil
}
