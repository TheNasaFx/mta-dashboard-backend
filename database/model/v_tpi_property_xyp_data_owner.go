package model

import "database/sql"

type PropertyOwner struct {
    ID                     sql.NullInt64   `db:"ID"`
    REG_NUM                sql.NullString  `db:"REG_NUM"`
    LAST_NAME              sql.NullString  `db:"LAST_NAME"`
    FIRST_NAME             sql.NullString  `db:"FIRST_NAME"`
    REGISTERED_DATE        sql.NullString  `db:"REGISTERED_DATE"`
    PROPERTY_NUMBER        sql.NullString  `db:"PROPERTY_NUMBER"`
    PROPERTY_SIZE          sql.NullString  `db:"PROPERTY_SIZE"`
    PROPERTY_VALUE         sql.NullString  `db:"PROPERTY_VALUE"`
    FULL_ADDRESS           sql.NullString  `db:"FULL_ADDRESS"`
    BRANCH_ID              sql.NullInt64   `db:"BRANCH_ID"`
    SUB_BRANCH_ID          sql.NullInt64   `db:"SUB_BRANCH_ID"`
    REF_PURPOSE_ID         sql.NullInt64   `db:"REF_PURPOSE_ID"`
    REF_SUB_PURPOSE_ID     sql.NullInt64   `db:"REF_SUB_PURPOSE_ID"`
    PROPERTY_TYPE          sql.NullString  `db:"PROPERTY_TYPE"`
    REF_SERVICE_TYPE_ID    sql.NullInt64   `db:"REF_SERVICE_TYPE_ID"`
    RE_OWNER_STATUS_ID     sql.NullInt64   `db:"RE_OWNER_STATUS_ID"`
    CREATED_BY             sql.NullString  `db:"CREATED_BY"`
    CREATED_DATE           sql.NullString  `db:"CREATED_DATE"`
    UPDATED_BY             sql.NullString  `db:"UPDATED_BY"`
    UPDATED_DATE           sql.NullString  `db:"UPDATED_DATE"`
    ACTIVE_FLAG            sql.NullString  `db:"ACTIVE_FLAG"`
    PRIMARY_ID             sql.NullInt64   `db:"PRIMARY_ID"`
    VERSION                sql.NullInt64   `db:"VERSION"`
    STATUS                 sql.NullString  `db:"STATUS"`
    ACCESS_LEVEL           sql.NullString  `db:"ACCESS_LEVEL"`
    ACTION_FLAG            sql.NullString  `db:"ACTION_FLAG"`
    REGISTER_STATUS_ID     sql.NullInt64   `db:"REGISTER_STATUS_ID"`
    IS_TAIS                sql.NullString  `db:"IS_TAIS"`
    ENT_ID                 sql.NullInt64   `db:"ENT_ID"`
    TAXPAYER_BRANCH_ID     sql.NullInt64   `db:"TAXPAYER_BRANCH_ID"`
    TAXPAYER_SUB_BRANCH_ID sql.NullInt64   `db:"TAXPAYER_SUB_BRANCH_ID"`
    TAXPAYER_BRANCH_NAME   sql.NullString  `db:"TAXPAYER_BRANCH_NAME"`
    TAXPAYER_SUB_BRANCH_NAME sql.NullString `db:"TAXPAYER_SUB_BRANCH_NAME"`
    TAXPAYER_REGISTERED    sql.NullString  `db:"TAXPAYER_REGISTERED"`
} 