package repository

import (
    "database/sql"
    "dashboard-backend/database/model"
)

func GetAccountGeneralYears(db *sql.DB) ([]model.AccountGeneralYear, error) {
    rows, err := db.Query("SELECT * FROM GPS.V_ACCOUNT_GENERAL_YEAR")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []model.AccountGeneralYear
    for rows.Next() {
        var agy model.AccountGeneralYear
        err := rows.Scan(
            &agy.PIN, &agy.ENTITY_NAME, &agy.TAX_DTYPE_CODE, &agy.TAX_DTYPE_NAME, &agy.ENT_ID, &agy.YEAR, &agy.BRANCH_NAME, &agy.C2_CREDIT, &agy.C2_DEBIT, &agy.PAYABLE_DEBIT, &agy.PAYABLE_CREDIT, &agy.PAYABLE_CONFIG, &agy.PERIOD_TYPE, &agy.TAX_TYPE_NAME, &agy.ACCOUNT_ID, &agy.C1_CREDIT, &agy.C1_DEBIT, &agy.TAX_TYPE_CODE,
        )
        if err != nil {
            return nil, err
        }
        results = append(results, agy)
    }
    return results, nil
} 