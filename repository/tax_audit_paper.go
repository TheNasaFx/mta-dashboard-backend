package repository

import (
    "database/sql"
    "dashboard-backend/database/model"
)

func GetTaxAuditPapers(db *sql.DB) ([]model.TaxAuditPaper, error) {
    rows, err := db.Query("SELECT * FROM GPS.V_TAX_AUDIT_PAPER")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []model.TaxAuditPaper
    for rows.Next() {
        var tap model.TaxAuditPaper
        err := rows.Scan(
            &tap.TAPR_SID, &tap.TAPR_ACC_SID, &tap.TAPR_MRCH_SID, &tap.TAPR_MRCH_OFF_CODE, &tap.TAPR_CODE, &tap.TAPR_DATE, &tap.TAPR_SDATE, &tap.TAPR_PAPER_TYPE,
        )
        if err != nil {
            return nil, err
        }
        results = append(results, tap)
    }
    return results, nil
} 