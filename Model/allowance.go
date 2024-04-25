package Model

import (
	"database/sql"
	"errors"
	"fmt"
)

// Custom errors
var (
	ErrNotFound     = errors.New("record not found")
	ErrInvalidRange = errors.New("new value is not within the allowed range")
)

// Allowance represents a record in the allowance table
type Allowance struct {
	ID            int
	AllowanceType string
	Min           int
	Max           int
	SetValue      int
}

// updateAllowanceSetValues updates the set value of an allowance record based on the given allowance type
func updateAllowanceSetValues(db *sql.DB, allowanceType string, newValue int) error {
	query := fmt.Sprintf(`SELECT id, min, max, set_value FROM allowance WHERE allowance_type = '%s' LIMIT 1`, allowanceType)

	var allowance Allowance
	err := db.QueryRow(query).Scan(&allowance.ID, &allowance.Min, &allowance.Max, &allowance.SetValue)
	if err != nil {
		return fmt.Errorf("failed to retrieve allowance record: %v", err)
	}

	if newValue < allowance.Min || newValue > allowance.Max {
		return ErrInvalidRange
	}

	updateQuery := `UPDATE allowance SET set_value = $1 WHERE id = $2`

	_, err = db.Exec(updateQuery, newValue, allowance.ID)
	if err != nil {
		return fmt.Errorf("failed to update allowance set values: %v", err)
	}

	fmt.Println("Allowance set values updated successfully")

	return nil
}

// UpdateKReceiptSetValues updates the set value of a k-receipt allowance record
func UpdateKReceiptSetValues(db *sql.DB, newValue int) error {
	return updateAllowanceSetValues(db, "k-receipt", newValue)
}

// UpdatePersonalSetValues updates the set value of a personal allowance record
func UpdatePersonalSetValues(db *sql.DB, newValue int) error {
	return updateAllowanceSetValues(db, "personal", newValue)
}

// GetAllAllowances retrieves all allowances from the database
func GetAllAllowances(db *sql.DB) ([]Allowance, error) {
	query := `SELECT id, allowance_type, min, max, set_value FROM allowance`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch allowances: %v", err)
	}
	defer rows.Close()

	var allowancesDB []Allowance
	for rows.Next() {
		var allowance Allowance
		err := rows.Scan(&allowance.ID, &allowance.AllowanceType, &allowance.Min, &allowance.Max, &allowance.SetValue)
		if err != nil {
			return nil, fmt.Errorf("failed to scan allowance row: %v", err)
		}
		allowancesDB = append(allowancesDB, allowance)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over allowance rows: %v", err)
	}

	return allowancesDB, nil
}
