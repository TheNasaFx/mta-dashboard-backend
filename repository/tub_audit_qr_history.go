package repository

import (
	"dashboard-backend/database/model"
	"database/sql"
)

func GetTubAuditQrHistories(db *sql.DB) ([]model.TubAuditQrHistory, error) {
	rows, err := db.Query("SELECT * FROM GPS.V_TAX_AUDIT_QR_HISTORY")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.TubAuditQrHistory
	for rows.Next() {
		var taqh model.TubAuditQrHistory
		err := rows.Scan(
			&taqh.ID, &taqh.SCANNED_DATE, &taqh.LONGITUDE, &taqh.LATITUDE, &taqh.QR_DATA, &taqh.WORKER_ID, &taqh.WORKER_USERNAME, &taqh.GOAL, &taqh.DESCRIPTION, &taqh.REGNO, &taqh.FINISHED_DATE,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, taqh)
	}
	return results, nil
}

func GetTubAuditQrHistoriesByRegno(db *sql.DB, regno string) ([]model.TubAuditQrHistory, error) {
	rows, err := db.Query("SELECT * FROM GPS.V_TAX_AUDIT_QR_HISTORY WHERE REGNO = :1", regno)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.TubAuditQrHistory
	for rows.Next() {
		var taqh model.TubAuditQrHistory
		err := rows.Scan(
			&taqh.ID, &taqh.SCANNED_DATE, &taqh.LONGITUDE, &taqh.LATITUDE, &taqh.QR_DATA, &taqh.WORKER_ID, &taqh.WORKER_USERNAME, &taqh.GOAL, &taqh.DESCRIPTION, &taqh.REGNO, &taqh.FINISHED_DATE,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, taqh)
	}
	return results, nil
}
