package model

import "database/sql"

type Payment struct {
    ID                  sql.NullInt64   `db:"ID"`
    SRC_ACCOUNT_ID      sql.NullInt64   `db:"SRC_ACCOUNT_ID"`
    DEST_ACCOUNT_ID     sql.NullInt64   `db:"DEST_ACCOUNT_ID"`
    CURRENCY_RATE       sql.NullString  `db:"CURRENCY_RATE"`
    FEE                 sql.NullString  `db:"FEE"`
    OWNER_ID            sql.NullInt64   `db:"OWNER_ID"`
    INVOICE_ID          sql.NullInt64   `db:"INVOICE_ID"`
    BANK_TRAN_NO        sql.NullString  `db:"BANK_TRAN_NO"`
    PAID_DATE           sql.NullString  `db:"PAID_DATE"`
    PAY_TYPE_ID         sql.NullInt64   `db:"PAY_TYPE_ID"`
    TRAN_TYPE           sql.NullString  `db:"TRAN_TYPE"`
    SETTLEMENT_ID       sql.NullInt64   `db:"SETTLEMENT_ID"`
    STATUS              sql.NullString  `db:"STATUS"`
    CREATED_BY          sql.NullString  `db:"CREATED_BY"`
    CREATED_DATE        sql.NullString  `db:"CREATED_DATE"`
    UPDATED_BY          sql.NullString  `db:"UPDATED_BY"`
    UPDATED_DATE        sql.NullString  `db:"UPDATED_DATE"`
    DESCRIPTION         sql.NullString  `db:"DESCRIPTION"`
    VERSION             sql.NullInt64   `db:"VERSION"`
    ACCESS_LEVEL        sql.NullString  `db:"ACCESS_LEVEL"`
    ACTIVE_FLAG         sql.NullString  `db:"ACTIVE_FLAG"`
    PRIMARY_ID          sql.NullInt64   `db:"PRIMARY_ID"`
    ACTION_FLAG         sql.NullString  `db:"ACTION_FLAG"`
    SRC_ACCOUNT_TYPE    sql.NullString  `db:"SRC_ACCOUNT_TYPE"`
    INV_TYPE            sql.NullString  `db:"INV_TYPE"`
    BANK_ID             sql.NullInt64   `db:"BANK_ID"`
    OPERATOR_ID         sql.NullInt64   `db:"OPERATOR_ID"`
    STATEMENT_STATUS    sql.NullString  `db:"STATEMENT_STATUS"`
    ADR_CONTACT_ID      sql.NullInt64   `db:"ADR_CONTACT_ID"`
    RECORD_SOURCE       sql.NullString  `db:"RECORD_SOURCE"`
    SETTLEMENT_NO       sql.NullString  `db:"SETTLEMENT_NO"`
    TMP_TRAN_ID         sql.NullString  `db:"TMP_TRAN_ID"`
    TMP_TAXACT_DLN      sql.NullString  `db:"TMP_TAXACT_DLN"`
    SUB_BUDGET_ID       sql.NullInt64   `db:"SUB_BUDGET_ID"`
    STATE_STATEMENT_ID  sql.NullInt64   `db:"STATE_STATEMENT_ID"`
    STATE_SETTLEMENT_DATE sql.NullString `db:"STATE_SETTLEMENT_DATE"`
    AMOUNT              sql.NullString  `db:"AMOUNT"`
    INVOICE_NO          sql.NullString  `db:"INVOICE_NO"`
    PAY_UUID            sql.NullString  `db:"PAY_UUID"`
    ACT_ACCOUNT_ID      sql.NullInt64   `db:"ACT_ACCOUNT_ID"`
    TAX_TYPE_ID         sql.NullInt64   `db:"TAX_TYPE_ID"`
    TAX_DTYPE_ID        sql.NullInt64   `db:"TAX_DTYPE_ID"`
    BRANCH_ID           sql.NullInt64   `db:"BRANCH_ID"`
    SUB_BRANCH_ID       sql.NullInt64   `db:"SUB_BRANCH_ID"`
    FIN_TRAN_NO         sql.NullString  `db:"FIN_TRAN_NO"`
    ACCOUNT_NO          sql.NullString  `db:"ACCOUNT_NO"`
} 