package repository

import (
    "database/sql"
    "dashboard-backend/database/model"
)

func GetPropertyOwners(db *sql.DB) ([]model.PropertyOwner, error) {
    rows, err := db.Query("SELECT * FROM GPS.V_TPI_PROPERTY_XYP_DATA_OWNER")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []model.PropertyOwner
    for rows.Next() {
        var po model.PropertyOwner
        err := rows.Scan(
            &po.ID, &po.REG_NUM, &po.LAST_NAME, &po.FIRST_NAME, &po.REGISTERED_DATE, &po.PROPERTY_NUMBER,
            &po.PROPERTY_SIZE, &po.PROPERTY_VALUE, &po.FULL_ADDRESS, &po.BRANCH_ID, &po.SUB_BRANCH_ID,
            &po.REF_PURPOSE_ID, &po.REF_SUB_PURPOSE_ID, &po.PROPERTY_TYPE, &po.REF_SERVICE_TYPE_ID,
            &po.RE_OWNER_STATUS_ID, &po.CREATED_BY, &po.CREATED_DATE, &po.UPDATED_BY, &po.UPDATED_DATE,
            &po.ACTIVE_FLAG, &po.PRIMARY_ID, &po.VERSION, &po.STATUS, &po.ACCESS_LEVEL, &po.ACTION_FLAG,
            &po.REGISTER_STATUS_ID, &po.IS_TAIS, &po.ENT_ID, &po.TAXPAYER_BRANCH_ID, &po.TAXPAYER_SUB_BRANCH_ID,
            &po.TAXPAYER_BRANCH_NAME, &po.TAXPAYER_SUB_BRANCH_NAME, &po.TAXPAYER_REGISTERED,
        )
        if err != nil {
            return nil, err
        }
        results = append(results, po)
    }
    return results, nil
} 