package model

import "database/sql"

type TubHrmWorkPerformance struct {
    CODE                        sql.NullString  `db:"CODE"`
    ID                          sql.NullInt64   `db:"ID"`
    WORKER_ID                   sql.NullInt64   `db:"WORKER_ID"`
    REF_WORK_PERFORMANCE_TYPE_ID sql.NullInt64  `db:"REF_WORK_PERFORMANCE_TYPE_ID"`
    WORK_SOLUTION               sql.NullString  `db:"WORK_SOLUTION"`
    NOTED_DATE                  sql.NullString  `db:"NOTED_DATE"`
    STATUS                      sql.NullString  `db:"STATUS"`
    CREATED_BY                  sql.NullString  `db:"CREATED_BY"`
    CREATED_DATE                sql.NullString  `db:"CREATED_DATE"`
    UPDATED_BY                  sql.NullString  `db:"UPDATED_BY"`
    UPDATED_DATE                sql.NullString  `db:"UPDATED_DATE"`
    VERSION                     sql.NullInt64   `db:"VERSION"`
    ACCESS_LEVEL                sql.NullString  `db:"ACCESS_LEVEL"`
    ACTIVE_FLAG                 sql.NullString  `db:"ACTIVE_FLAG"`
    PRIMARY_ID                  sql.NullInt64   `db:"PRIMARY_ID"`
    ACTION_FLAG                 sql.NullString  `db:"ACTION_FLAG"`
    WORK_NAME                   sql.NullString  `db:"WORK_NAME"`
} 