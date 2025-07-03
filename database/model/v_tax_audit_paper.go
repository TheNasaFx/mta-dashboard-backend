package model

import "database/sql"

type TaxAuditPaper struct {
    TAPR_SID         sql.NullInt64   `db:"TAPR_SID"`
    TAPR_ACC_SID     sql.NullInt64   `db:"TAPR_ACC_SID"`
    TAPR_MRCH_SID    sql.NullInt64   `db:"TAPR_MRCH_SID"`
    TAPR_MRCH_OFF_CODE sql.NullString `db:"TAPR_MRCH_OFF_CODE"`
    TAPR_CODE        sql.NullString  `db:"TAPR_CODE"`
    TAPR_DATE        sql.NullString  `db:"TAPR_DATE"`
    TAPR_SDATE       sql.NullString  `db:"TAPR_SDATE"`
    TAPR_PAPER_TYPE  sql.NullString  `db:"TAPR_PAPER_TYPE"`
} 