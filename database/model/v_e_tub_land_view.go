package model

import "database/sql"

// LandView represents a row from GPS.V_E_TUB_LAND_VIEW
type LandView struct {
	ID                         int             `db:"ID"`
	ENT_ID                     int             `db:"ENT_ID"`
	PIN                        sql.NullString  `db:"PIN"`
	NAME                       sql.NullString  `db:"NAME"`
	OBJ_TYPE                   sql.NullString  `db:"OBJ_TYPE"`
	CERTIFICATE_NO             sql.NullString  `db:"CERTIFICATE_NO"`
	CERTIFICATE_DATE           sql.NullString  `db:"CERTIFICATE_DATE"`
	START_DATE                 sql.NullString  `db:"START_DATE"`
	END_DATE                   sql.NullString  `db:"END_DATE"`
	BRANCH_ID                  sql.NullInt64   `db:"BRANCH_ID"`
	DESCRIPTION                sql.NullString  `db:"DESCRIPTION"`
	STATUS                     sql.NullString  `db:"STATUS"`
	ACTIVE_FLAG                sql.NullString  `db:"ACTIVE_FLAG"`
	SUB_BRANCH_ID              sql.NullInt64   `db:"SUB_BRANCH_ID"`
	LAND_ID                    sql.NullInt64   `db:"LAND_ID"`
	PARCEL_ID                  sql.NullInt64   `db:"PARCEL_ID"`
	REGISTER_NO                sql.NullString  `db:"REGISTER_NO"`
	AREA_M2                    sql.NullFloat64 `db:"AREA_M2"`
	AU1_CODE                   sql.NullString  `db:"AU1_CODE"`
	AU1_NAME                   sql.NullString  `db:"AU1_NAME"`
	AU2_CODE                   sql.NullString  `db:"AU2_CODE"`
	AU2_NAME                   sql.NullString  `db:"AU2_NAME"`
	AU3_CODE                   sql.NullString  `db:"AU3_CODE"`
	AU3_NAME                   sql.NullString  `db:"AU3_NAME"`
	ADDRESS_STREETNAME         sql.NullString  `db:"ADDRESS_STREETNAME"`
	ADDRESS_KHASHAA            sql.NullString  `db:"ADDRESS_KHASHAA"`
	REF_LAND_DEDICATION_ID     sql.NullInt64   `db:"REF_LAND_DEDICATION_ID"`
	REF_LAND_DECISION_LEVEL_ID sql.NullInt64   `db:"REF_LAND_DECISION_LEVEL_ID"`
	REF_LAND_TYPE_ID           sql.NullInt64   `db:"REF_LAND_TYPE_ID"`
	OBJ_ID                     sql.NullInt64   `db:"OBJ_ID"`
	DECISION_DATE              sql.NullString  `db:"DECISION_DATE"`
	DECISION_NO                sql.NullString  `db:"DECISION_NO"`
	COORD_Y                    sql.NullString  `db:"COORD_Y"`
	COORD_X                    sql.NullString  `db:"COORD_X"`
}
