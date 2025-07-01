package repository

import (
	"dashboard-backend/database"
	"dashboard-backend/database/model"
	"database/sql"
	"fmt"
)

// GetOrgList returns a list of organizations with optional filters and pagination
func GetOrgList(name, code, status string, pageSize, pageNumber int) ([]model.Org, error) {
	if database.DB == nil {
		database.MustConnect()
	}

	// Oracle 11g-compatible pagination using ROWNUM
	query := `SELECT * FROM (
		SELECT a.*, ROWNUM rnum FROM (
			SELECT ID, NAME, OFFICE_CODE, REGNO, KHO_CODE, BUILD_FLOOR, ADDRESS, LNG, LAT
			FROM GPS.PAY_CENTER
			WHERE 1=1
			ORDER BY ID DESC
		) a WHERE ROWNUM <= :maxRow
	) WHERE rnum > :minRow`

	minRow := (pageNumber - 1) * pageSize
	maxRow := pageNumber * pageSize

	rows, err := database.DB.Query(query, maxRow, minRow)
	if err != nil {
		return nil, fmt.Errorf("DB query error: %w", err)
	}
	defer rows.Close()

	var orgs []model.Org
	for rows.Next() {
		var org model.Org
		var rnum int
		err := rows.Scan(
			&org.ID,
			&org.Name,
			&org.OfficeCode,
			&org.Regno,
			&org.KhoCode,
			&org.BuildFloor,
			&org.Address,
			&org.Lng,
			&org.Lat,
			&rnum,
		)
		if err != nil {
			return nil, fmt.Errorf("Scan error: %w", err)
		}
		orgs = append(orgs, org)
	}
	return orgs, nil
}

// FindOrgByID returns a single organization by ID
func FindOrgByID(id uint) (*model.Org, error) {
	if database.DB == nil {
		database.MustConnect()
	}

	query := `SELECT ID, NAME, OFFICE_CODE, REGNO, KHO_CODE, BUILD_FLOOR, ADDRESS, LNG, LAT FROM GPS.PAY_CENTER WHERE ID = :id`
	row := database.DB.QueryRow(query, id)

	var org model.Org
	err := row.Scan(
		&org.ID,
		&org.Name,
		&org.OfficeCode,
		&org.Regno,
		&org.KhoCode,
		&org.BuildFloor,
		&org.Address,
		&org.Lng,
		&org.Lat,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Scan error: %w", err)
	}
	return &org, nil
}

type OrgInput struct {
	Name       string  `json:"name"`
	OfficeCode string  `json:"office_code"`
	Regno      string  `json:"regno"`
	KhoCode    string  `json:"kho_code"`
	BuildFloor string  `json:"build_floor"`
	Address    string  `json:"address"`
	Lng        float64 `json:"lng"`
	Lat        float64 `json:"lat"`
}

func CreateOrg(input OrgInput) (*model.Org, error) {
	if database.DB == nil {
		database.MustConnect()
	}
	query := `INSERT INTO GPS.PAY_CENTER (NAME, OFFICE_CODE, REGNO, KHO_CODE, BUILD_FLOOR, ADDRESS, LNG, LAT)
		VALUES (:1, :2, :3, :4, :5, :6, :7, :8) RETURNING ID INTO :9`
	var id uint
	_, err := database.DB.Exec(query,
		input.Name, input.OfficeCode, input.Regno, input.KhoCode, input.BuildFloor, input.Address, input.Lng, input.Lat, &id)
	if err != nil {
		return nil, err
	}
	return &model.Org{
		ID:         id,
		Name:       input.Name,
		OfficeCode: input.OfficeCode,
		Regno:      input.Regno,
		KhoCode:    input.KhoCode,
		BuildFloor: input.BuildFloor,
		Address:    input.Address,
		Lng:        input.Lng,
		Lat:        input.Lat,
	}, nil
}

func UpdateOrg(id uint, input OrgInput) (*model.Org, error) {
	if database.DB == nil {
		database.MustConnect()
	}
	query := `UPDATE GPS.PAY_CENTER SET NAME=:1, OFFICE_CODE=:2, REGNO=:3, KHO_CODE=:4, BUILD_FLOOR=:5, ADDRESS=:6, LNG=:7, LAT=:8 WHERE ID=:9`
	_, err := database.DB.Exec(query,
		input.Name, input.OfficeCode, input.Regno, input.KhoCode, input.BuildFloor, input.Address, input.Lng, input.Lat, id)
	if err != nil {
		return nil, err
	}
	return FindOrgByID(id)
}

func DeleteOrg(id uint) error {
	if database.DB == nil {
		database.MustConnect()
	}
	query := `DELETE FROM GPS.PAY_CENTER WHERE ID=:1`
	_, err := database.DB.Exec(query, id)
	return err
}
