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

// UpdateAllowanceSetValues updates set_value in allowance table where allowance_type is 'k-recipe'
func UpdateAllowanceSetValues(db *sql.DB, newValue int) error {
	// Retrieve the current value, min value, and max value from the database
	query := `SELECT id, min, max, set_value FROM allowance WHERE allowance_type = 'k-recipe' LIMIT 1`

	var allowance Allowance
	err := db.QueryRow(query).Scan(&allowance.ID, &allowance.Min, &allowance.Max, &allowance.SetValue)
	if err != nil {
		return fmt.Errorf("failed to retrieve allowance record: %v", err)
	}

	// Validate the newValue against min and max values
	if newValue < allowance.Min || newValue > allowance.Max {
		return errors.New("newValue is not within the allowed range")
	}

	// Prepare the SQL query
	updateQuery := `UPDATE allowance SET set_value = $1 WHERE id = $2`

	// Execute the query
	_, err = db.Exec(updateQuery, newValue, allowance.ID)
	if err != nil {
		return fmt.Errorf("failed to update allowance set values: %v", err)
	}

	fmt.Println("Allowance set values updated successfully")

	return nil
}
