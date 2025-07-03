package repository

import (
    "database/sql"
    "dashboard-backend/database/model"
)

func GetTubReportData(db *sql.DB) ([]model.TubReportData, error) {
    rows, err := db.Query("SELECT * FROM GPS.V_E_TUB_REPORT_DATA")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []model.TubReportData
    for rows.Next() {
        var trd model.TubReportData
        err := rows.Scan(
            &trd.PIN, &trd.MAINTYPE_CODE, &trd.OFFICE_CODE, &trd.TAX_REPORT_CODE, &trd.FREQUENCY, &trd.TAX_YEAR, &trd.TAX_PERIOD, &trd.WORKFLOW_STATUS_ID, &trd.CREATED_DATE, &trd.UPDATED_DATE, &trd.RECEIPT_DATE, &trd.IS_ACTIVE, &trd.DONE_DATE, &trd.SUBMITTED_DATE, &trd.ENT_ID, &trd.BRANCH_ID,
        )
        if err != nil {
            return nil, err
        }
        results = append(results, trd)
    }
    return results, nil
} 