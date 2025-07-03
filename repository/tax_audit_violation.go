package repository

import (
    "database/sql"
    "dashboard-backend/database/model"
)

func GetTaxAuditViolations(db *sql.DB) ([]model.TaxAuditViolation, error) {
    rows, err := db.Query("SELECT * FROM GPS.V_TAX_AUDIT_VIOLATION")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []model.TaxAuditViolation
    for rows.Next() {
        var tav model.TaxAuditViolation
        err := rows.Scan(
            &tav.TAPR_SID, &tav.TAV_VIOLATION_CODE, &tav.TAV_NOTICE_DATE, &tav.TAV_ELIMINATE_VIOLATION_DAY, &tav.TAV_STATUS,
        )
        if err != nil {
            return nil, err
        }
        results = append(results, tav)
    }
    return results, nil
} 