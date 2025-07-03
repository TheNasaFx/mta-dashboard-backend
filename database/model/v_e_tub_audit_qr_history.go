package model

import "database/sql"

type TubAuditQrHistory struct {
    ID              sql.NullInt64   `db:"ID"`
    SCANNED_DATE    sql.NullString  `db:"SCANNED_DATE"`
    LONGITUDE       sql.NullString  `db:"LONGITUDE"`
    LATITUDE        sql.NullString  `db:"LATITUDE"`
    QR_DATA         sql.NullString  `db:"QR_DATA"`
    WORKER_ID       sql.NullInt64   `db:"WORKER_ID"`
    WORKER_USERNAME sql.NullString  `db:"WORKER_USERNAME"`
    GOAL            sql.NullString  `db:"GOAL"`
    DESCRIPTION     sql.NullString  `db:"DESCRIPTION"`
    REGNO           sql.NullString  `db:"REGNO"`
    FINISHED_DATE   sql.NullString  `db:"FINISHED_DATE"`
} 