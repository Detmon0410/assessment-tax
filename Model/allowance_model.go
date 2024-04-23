// model_allowance.go

package Model

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Allowance struct {
	ID            int    `json:"id"`
	AllowanceType string `json:"allowance_type"`
	Min           int    `json:"min"`
	Max           int    `json:"max"`
	SetValue      int    `json:"set_value"`
}

func QueryAllowancesFromDB() ([]Allowance, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("connect to database error: %v", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, allowance_type, min, max, set_value FROM allowance")
	if err != nil {
		return nil, fmt.Errorf("can't prepare query statement: %v", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("can't query allowances: %v", err)
	}

	var allowancesDB []Allowance
	for rows.Next() {
		var allowance Allowance
		err := rows.Scan(&allowance.ID, &allowance.AllowanceType, &allowance.Min, &allowance.Max, &allowance.SetValue)
		if err != nil {
			return nil, fmt.Errorf("can't scan row into variable: %v", err)
		}
		allowancesDB = append(allowancesDB, allowance)
	}
	fmt.Println("Allowances from DB:", allowancesDB)
	return allowancesDB, nil
}
