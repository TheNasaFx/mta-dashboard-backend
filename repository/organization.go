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

	query := `SELECT id, name, code, type, address, phone, email, status, parent_id, created_at, updated_at, deleted_at FROM pay_center.org WHERE 1=1`
	args := []interface{}{}

	if name != "" {
		query += " AND LOWER(name) LIKE LOWER(?)"
		args = append(args, "%"+name+"%")
	}
	if code != "" {
		query += " AND LOWER(code) LIKE LOWER(?)"
		args = append(args, "%"+code+"%")
	}
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	// Paging
	offset := (pageNumber - 1) * pageSize
	query += " ORDER BY created_at DESC OFFSET ? ROWS FETCH NEXT ? ROWS ONLY"
	args = append(args, offset, pageSize)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("DB query error: %w", err)
	}
	defer rows.Close()

	var orgs []model.Org
	for rows.Next() {
		var org model.Org
		err := rows.Scan(
			&org.ID, &org.Name,
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

	query := `SELECT ID, NAME FROM GPS.PAY_CENTER`
	row := database.DB.QueryRow(query, id)

	var org model.Org
	err := row.Scan(
		&org.ID, &org.Name,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Scan error: %w", err)
	}
	return &org, nil
}
