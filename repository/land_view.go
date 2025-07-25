package repository

import (
	"dashboard-backend/database/model"
	"database/sql"
)

func GetLandViews(db *sql.DB) ([]model.LandView, error) {
	rows, err := db.Query("SELECT * FROM GPS.V_E_TUB_LAND_VIEW")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.LandView
	for rows.Next() {
		var lv model.LandView
		err := rows.Scan(
			&lv.ID, &lv.ENT_ID, &lv.PIN, &lv.NAME, &lv.OBJ_TYPE, &lv.CERTIFICATE_NO, &lv.CERTIFICATE_DATE,
			&lv.START_DATE, &lv.END_DATE, &lv.BRANCH_ID, &lv.DESCRIPTION, &lv.STATUS, &lv.ACTIVE_FLAG,
			&lv.SUB_BRANCH_ID, &lv.LAND_ID, &lv.PARCEL_ID, &lv.REGISTER_NO, &lv.AREA_M2,
			&lv.AU1_CODE, &lv.AU1_NAME, &lv.AU2_CODE, &lv.AU2_NAME, &lv.AU3_CODE, &lv.AU3_NAME,
			&lv.ADDRESS_STREETNAME, &lv.ADDRESS_KHASHAA, &lv.REF_LAND_DEDICATION_ID, &lv.REF_LAND_DECISION_LEVEL_ID,
			&lv.REF_LAND_TYPE_ID, &lv.OBJ_ID, &lv.DECISION_DATE, &lv.DECISION_NO, &lv.COORD_Y, &lv.COORD_X,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, lv)
	}
	return results, nil
}
