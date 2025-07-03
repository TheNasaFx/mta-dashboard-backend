package repository

import (
    "database/sql"
    "dashboard-backend/database/model"
)

func GetPayments(db *sql.DB) ([]model.Payment, error) {
    rows, err := db.Query("SELECT * FROM GPS.V_E_TUB_PAYMENTS")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []model.Payment
    for rows.Next() {
        var p model.Payment
        err := rows.Scan(
            &p.ID, &p.SRC_ACCOUNT_ID, &p.DEST_ACCOUNT_ID, &p.CURRENCY_RATE, &p.FEE, &p.OWNER_ID, &p.INVOICE_ID, &p.BANK_TRAN_NO, &p.PAID_DATE, &p.PAY_TYPE_ID, &p.TRAN_TYPE, &p.SETTLEMENT_ID, &p.STATUS, &p.CREATED_BY, &p.CREATED_DATE, &p.UPDATED_BY, &p.UPDATED_DATE, &p.DESCRIPTION, &p.VERSION, &p.ACCESS_LEVEL, &p.ACTIVE_FLAG, &p.PRIMARY_ID, &p.ACTION_FLAG, &p.SRC_ACCOUNT_TYPE, &p.INV_TYPE, &p.BANK_ID, &p.OPERATOR_ID, &p.STATEMENT_STATUS, &p.ADR_CONTACT_ID, &p.RECORD_SOURCE, &p.SETTLEMENT_NO, &p.TMP_TRAN_ID, &p.TMP_TAXACT_DLN, &p.SUB_BUDGET_ID, &p.STATE_STATEMENT_ID, &p.STATE_SETTLEMENT_DATE, &p.AMOUNT, &p.INVOICE_NO, &p.PAY_UUID, &p.ACT_ACCOUNT_ID, &p.TAX_TYPE_ID, &p.TAX_DTYPE_ID, &p.BRANCH_ID, &p.SUB_BRANCH_ID, &p.FIN_TRAN_NO, &p.ACCOUNT_NO,
        )
        if err != nil {
            return nil, err
        }
        results = append(results, p)
    }
    return results, nil
} 