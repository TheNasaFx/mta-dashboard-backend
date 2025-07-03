package repository

import (
    "database/sql"
    "dashboard-backend/database/model"
)

func GetTubHrmWorkPerformances(db *sql.DB) ([]model.TubHrmWorkPerformance, error) {
    rows, err := db.Query("SELECT * FROM GPS.V_E_TUB_HRM_WORK_PERFORMANCE")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []model.TubHrmWorkPerformance
    for rows.Next() {
        var twhp model.TubHrmWorkPerformance
        err := rows.Scan(
            &twhp.CODE, &twhp.ID, &twhp.WORKER_ID, &twhp.REF_WORK_PERFORMANCE_TYPE_ID, &twhp.WORK_SOLUTION, &twhp.NOTED_DATE, &twhp.STATUS, &twhp.CREATED_BY, &twhp.CREATED_DATE, &twhp.UPDATED_BY, &twhp.UPDATED_DATE, &twhp.VERSION, &twhp.ACCESS_LEVEL, &twhp.ACTIVE_FLAG, &twhp.PRIMARY_ID, &twhp.ACTION_FLAG, &twhp.WORK_NAME,
        )
        if err != nil {
            return nil, err
        }
        results = append(results, twhp)
    }
    return results, nil
} 