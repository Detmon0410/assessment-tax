package Model

import (
	"testing"

	// Include any other necessary packages for your specific testing environment
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAllAllowances(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define the expected rows returned by the query
	rows := sqlmock.NewRows([]string{"id", "allowance_type", "min", "max", "set_value"}).
		AddRow(1, "donation", 0, 100000, 100000).
		AddRow(3, "k-receipt", 0, 100000, 50000).
		AddRow(2, "personal", 10000, 100000, 60000)

	// Expect the query to be executed and return the defined rows
	mock.ExpectQuery("SELECT id, allowance_type, min, max, set_value FROM allowance").
		WillReturnRows(rows)

	// Call the function being tested
	allowances, err := GetAllAllowances(db)
	if err != nil {
		t.Errorf("error fetching allowances: %v", err)
		return
	}

	// Assert that the expected rows match the actual result
	assert.Equal(t, 3, len(allowances), "expected three allowances")
	assert.Equal(t, "donation", allowances[0].AllowanceType)
	assert.Equal(t, 100000, allowances[0].SetValue)
	assert.Equal(t, "k-receipt", allowances[1].AllowanceType)
	assert.Equal(t, 50000, allowances[1].SetValue)
	assert.Equal(t, "personal", allowances[2].AllowanceType)
	assert.Equal(t, 60000, allowances[2].SetValue)

	// Ensure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
