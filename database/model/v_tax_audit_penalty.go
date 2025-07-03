package model

import "database/sql"

type TaxAuditPenalty struct {
    TAPR_SID              sql.NullInt64   `db:"TAPR_SID"`
    TAP_REL_VIOLATION_CODE sql.NullString  `db:"TAP_REL_VIOLATION_CODE"`
    TAP_VIOLATION_CODE    sql.NullString  `db:"TAP_VIOLATION_CODE"`
    TAP_PENALTY_CODE      sql.NullString  `db:"TAP_PENALTY_CODE"`
    TAP_BNK_SID           sql.NullInt64   `db:"TAP_BNK_SID"`
    TAP_BNK_ACCOUNT       sql.NullString  `db:"TAP_BNK_ACCOUNT"`
    TAP_DUE_DATE          sql.NullString  `db:"TAP_DUE_DATE"`
    TAP_AMOUNT            sql.NullString  `db:"TAP_AMOUNT"`
} 