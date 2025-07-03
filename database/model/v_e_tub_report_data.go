package model

import "database/sql"

type TubReportData struct {
    PIN                sql.NullString  `db:"PIN"`
    MAINTYPE_CODE      sql.NullString  `db:"MAINTYPE_CODE"`
    OFFICE_CODE        sql.NullString  `db:"OFFICE_CODE"`
    TAX_REPORT_CODE    sql.NullString  `db:"TAX_REPORT_CODE"`
    FREQUENCY          sql.NullString  `db:"FREQUENCY"`
    TAX_YEAR           sql.NullInt64   `db:"TAX_YEAR"`
    TAX_PERIOD         sql.NullString  `db:"TAX_PERIOD"`
    WORKFLOW_STATUS_ID sql.NullInt64   `db:"WORKFLOW_STATUS_ID"`
    CREATED_DATE       sql.NullString  `db:"CREATED_DATE"`
    UPDATED_DATE       sql.NullString  `db:"UPDATED_DATE"`
    RECEIPT_DATE       sql.NullString  `db:"RECEIPT_DATE"`
    IS_ACTIVE          sql.NullString  `db:"IS_ACTIVE"`
    DONE_DATE          sql.NullString  `db:"DONE_DATE"`
    SUBMITTED_DATE     sql.NullString  `db:"SUBMITTED_DATE"`
    ENT_ID             sql.NullInt64   `db:"ENT_ID"`
    BRANCH_ID          sql.NullInt64   `db:"BRANCH_ID"`
} 