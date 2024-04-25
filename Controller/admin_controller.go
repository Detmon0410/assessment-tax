package Controller

import (
	"encoding/json"
	"net/http"

	"github.com/Detmon0410/assessment-tax/Model"
)

func UpdateAllowanceSetValuesHandler(w http.ResponseWriter, r *http.Request) {

	var requestData struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	setValue := int(requestData.Amount)

	db, err := Model.InitializeDB()
	if err != nil {
		http.Error(w, "Error initializing database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if err := Model.UpdateAllowanceSetValues(db, setValue); err != nil {
		http.Error(w, "Error updating allowance set values", http.StatusInternalServerError)
		return
	}

	responseData := struct {
		KReceipt float64 `json:"kReceipt"`
	}{
		KReceipt: requestData.Amount,
	}

	jsonResponse, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func GetAllAllowancesHandler(w http.ResponseWriter, r *http.Request) {

	db, err := Model.InitializeDB()
	if err != nil {
		http.Error(w, "Error initializing database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	allowances, err := Model.GetAllAllowances(db)
	if err != nil {
		http.Error(w, "Error fetching allowances", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(allowances)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
