// File: Model/allowance.go

package Model

import (
	"database/sql"
	"errors"
	"fmt"
)

// Allowance represents a record in the allowance table
type Allowance struct {
	ID            int
	AllowanceType string
	Min           int
	Max           int
	SetValue      int
}

func UpdateAllowanceSetValues(db *sql.DB, newValue int) error {

	query := `SELECT id, min, max, set_value FROM allowance WHERE allowance_type = 'k-receipt' LIMIT 1`

	var allowance Allowance
	err := db.QueryRow(query).Scan(&allowance.ID, &allowance.Min, &allowance.Max, &allowance.SetValue)
	if err != nil {
		return fmt.Errorf("failed to retrieve allowance record: %v", err)
	}

	if newValue < allowance.Min || newValue > allowance.Max {
		return errors.New("newValue is not within the allowed range")
	}

	updateQuery := `UPDATE allowance SET set_value = $1 WHERE id = $2`

	_, err = db.Exec(updateQuery, newValue, allowance.ID)
	if err != nil {
		return fmt.Errorf("failed to update allowance set values: %v", err)
	}

	fmt.Println("Allowance set values updated successfully")

	return nil
}

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
