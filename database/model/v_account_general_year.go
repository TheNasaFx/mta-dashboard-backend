package model

import "database/sql"

type AccountGeneralYear struct {
    PIN             sql.NullString  `db:"PIN"`
    ENTITY_NAME     sql.NullString  `db:"ENTITY_NAME"`
    TAX_DTYPE_CODE  sql.NullString  `db:"TAX_DTYPE_CODE"`
    TAX_DTYPE_NAME  sql.NullString  `db:"TAX_DTYPE_NAME"`
    ENT_ID          sql.NullInt64   `db:"ENT_ID"`
    YEAR            sql.NullInt64   `db:"YEAR"`
    BRANCH_NAME     sql.NullString  `db:"BRANCH_NAME"`
    C2_CREDIT       sql.NullString  `db:"C2_CREDIT"`
    C2_DEBIT        sql.NullString  `db:"C2_DEBIT"`
    PAYABLE_DEBIT   sql.NullString  `db:"PAYABLE_DEBIT"`
    PAYABLE_CREDIT  sql.NullString  `db:"PAYABLE_CREDIT"`
    PAYABLE_CONFIG  sql.NullString  `db:"PAYABLE_CONFIG"`
    PERIOD_TYPE     sql.NullString  `db:"PERIOD_TYPE"`
    TAX_TYPE_NAME   sql.NullString  `db:"TAX_TYPE_NAME"`
    ACCOUNT_ID      sql.NullInt64   `db:"ACCOUNT_ID"`
    C1_CREDIT       sql.NullString  `db:"C1_CREDIT"`
    C1_DEBIT        sql.NullString  `db:"C1_DEBIT"`
    TAX_TYPE_CODE   sql.NullString  `db:"TAX_TYPE_CODE"`
} 