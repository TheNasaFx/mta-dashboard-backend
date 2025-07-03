package model

import "database/sql"

type TaxAuditViolation struct {
    TAPR_SID                   sql.NullInt64   `db:"TAPR_SID"`
    TAV_VIOLATION_CODE         sql.NullString  `db:"TAV_VIOLATION_CODE"`
    TAV_NOTICE_DATE            sql.NullString  `db:"TAV_NOTICE_DATE"`
    TAV_ELIMINATE_VIOLATION_DAY sql.NullString `db:"TAV_ELIMINATE_VIOLATION_DAY"`
    TAV_STATUS                 sql.NullString  `db:"TAV_STATUS"`
} 