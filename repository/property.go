package repository

import (
	"dashboard-backend/database"
	"dashboard-backend/database/model"
	"database/sql"
	"fmt"
)

func GetPropertiesByPayCenterID(payCenterID uint) ([]model.Property, error) {
	if database.DB == nil {
		database.MustConnect()
	}
	query := `SELECT ID, PAY_CENTER_ID, UPDATED_DATE, PROPERTY_TYPE, OWNER_REGNO, PROPERTY_SIZE, RENT_AMOUNT
		FROM GPS.PAY_CENTER_PROPERTY
		WHERE PAY_CENTER_ID = :payCenterID`
	rows, err := database.DB.Query(query, payCenterID)
	if err != nil {
		return nil, fmt.Errorf("DB query error: %w", err)
	}
	defer rows.Close()

	var properties []model.Property
	for rows.Next() {
		var p model.Property
		err := rows.Scan(
			&p.ID,
			&p.PayCenterID,
			&p.UpdatedDate,
			&p.PropertyType,
			&p.OwnerRegno,
			&p.PropertySize,
			&p.RentAmount,
		)
		if err != nil {
			return nil, fmt.Errorf("Scan error: %w", err)
		}
		properties = append(properties, p)
	}
	return properties, nil
}

type PropertyInput struct {
	PayCenterID  uint     `json:"pay_center_id"`
	UpdatedDate  *string  `json:"updated_date"`
	PropertyType *string  `json:"property_type"`
	OwnerRegno   *string  `json:"owner_regno"`
	PropertySize *float64 `json:"property_size"`
	RentAmount   *float64 `json:"rent_amount"`
}

func CreateProperty(input PropertyInput) (*model.Property, error) {
	if database.DB == nil {
		database.MustConnect()
	}
	query := `INSERT INTO GPS.PAY_CENTER_PROPERTY (PAY_CENTER_ID, UPDATED_DATE, PROPERTY_TYPE, OWNER_REGNO, PROPERTY_SIZE, RENT_AMOUNT)
		VALUES (:1, :2, :3, :4, :5, :6) RETURNING ID INTO :7`
	var id uint
	_, err := database.DB.Exec(query,
		input.PayCenterID, input.UpdatedDate, input.PropertyType, input.OwnerRegno, input.PropertySize, input.RentAmount, &id)
	if err != nil {
		return nil, err
	}
	return &model.Property{
		ID:           id,
		PayCenterID:  input.PayCenterID,
		UpdatedDate:  sqlNullString(input.UpdatedDate),
		PropertyType: sqlNullString(input.PropertyType),
		OwnerRegno:   sqlNullString(input.OwnerRegno),
		PropertySize: sqlNullFloat64(input.PropertySize),
		RentAmount:   sqlNullFloat64(input.RentAmount),
	}, nil
}

func UpdateProperty(id uint, input PropertyInput) (*model.Property, error) {
	if database.DB == nil {
		database.MustConnect()
	}
	query := `UPDATE GPS.PAY_CENTER_PROPERTY SET PAY_CENTER_ID=:1, UPDATED_DATE=:2, PROPERTY_TYPE=:3, OWNER_REGNO=:4, PROPERTY_SIZE=:5, RENT_AMOUNT=:6 WHERE ID=:7`
	_, err := database.DB.Exec(query,
		input.PayCenterID, input.UpdatedDate, input.PropertyType, input.OwnerRegno, input.PropertySize, input.RentAmount, id)
	if err != nil {
		return nil, err
	}
	return GetPropertyByID(id)
}

func DeleteProperty(id uint) error {
	if database.DB == nil {
		database.MustConnect()
	}
	query := `DELETE FROM GPS.PAY_CENTER_PROPERTY WHERE ID=:1`
	_, err := database.DB.Exec(query, id)
	return err
}

func GetPropertyByID(id uint) (*model.Property, error) {
	if database.DB == nil {
		database.MustConnect()
	}
	query := `SELECT ID, PAY_CENTER_ID, UPDATED_DATE, PROPERTY_TYPE, OWNER_REGNO, PROPERTY_SIZE, RENT_AMOUNT FROM GPS.PAY_CENTER_PROPERTY WHERE ID = :1`
	row := database.DB.QueryRow(query, id)
	var p model.Property
	err := row.Scan(
		&p.ID,
		&p.PayCenterID,
		&p.UpdatedDate,
		&p.PropertyType,
		&p.OwnerRegno,
		&p.PropertySize,
		&p.RentAmount,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func sqlNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func sqlNullFloat64(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}
