package repository

import (
    "database/sql"
    "dashboard-backend/database/model"
)

func GetTaxAuditPenalties(db *sql.DB) ([]model.TaxAuditPenalty, error) {
    rows, err := db.Query("SELECT * FROM GPS.V_TAX_AUDIT_PENALTY")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []model.TaxAuditPenalty
    for rows.Next() {
        var tap model.TaxAuditPenalty
        err := rows.Scan(
            &tap.TAPR_SID, &tap.TAP_REL_VIOLATION_CODE, &tap.TAP_VIOLATION_CODE, &tap.TAP_PENALTY_CODE, &tap.TAP_BNK_SID, &tap.TAP_BNK_ACCOUNT, &tap.TAP_DUE_DATE, &tap.TAP_AMOUNT,
        )
        if err != nil {
            return nil, err
        }
        results = append(results, tap)
    }
    return results, nil
} 