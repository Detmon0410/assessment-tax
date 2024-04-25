// File: Controller/admin_controller.go

package Controller

import (
	"encoding/json"
	"net/http"

	"github.com/Detmon0410/assessment-tax/Model"
)

// UpdateAllowanceSetValuesHandler updates set_value in allowance table for the fixed allowance_type "k-recipe"
func UpdateAllowanceSetValuesHandler(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	var requestData struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	// Convert amount to integer (assuming set_value is an integer)
	setValue := int(requestData.Amount)

	// Initialize the database connection
	db, err := Model.InitializeDB()
	if err != nil {
		http.Error(w, "Error initializing database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Update the set value for 'k-recipe' allowance type
	if err := Model.UpdateAllowanceSetValues(db, setValue); err != nil {
		http.Error(w, "Error updating allowance set values", http.StatusInternalServerError)
		return
	}

	// Construct the response body
	responseData := struct {
		KReceipt float64 `json:"kReceipt"`
	}{
		KReceipt: requestData.Amount,
	}

	// Marshal response data into JSON
	jsonResponse, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	// Set Content-Type header and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// GetAllAllowancesHandler retrieves all records from the allowance table
func GetAllAllowancesHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize the database connection
	db, err := Model.InitializeDB()
	if err != nil {
		http.Error(w, "Error initializing database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Fetch all allowances from the database
	allowances, err := Model.GetAllAllowances(db)
	if err != nil {
		http.Error(w, "Error fetching allowances", http.StatusInternalServerError)
		return
	}

	// Marshal allowances into JSON
	jsonResponse, err := json.Marshal(allowances)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	// Set Content-Type header and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
