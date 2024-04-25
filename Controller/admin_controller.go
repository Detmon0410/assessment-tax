package Controller

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/Detmon0410/assessment-tax/Model"
	"github.com/joho/godotenv"
)

// Load environment variables from .env file
func init() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func basicAuth(w http.ResponseWriter, r *http.Request) (string, bool) {
	// Retrieve username and password from environment variables
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	// Check if username and password match the provided credentials
	username, password, ok := r.BasicAuth()
	if !ok || username != adminUsername || password != adminPassword {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return "", false
	}

	return username, true
}

func UpdateKReceiptHandler(w http.ResponseWriter, r *http.Request) {
	// Check basic authentication
	username, ok := basicAuth(w, r)
	if !ok {
		return
	}

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

	if err := Model.UpdateKReceiptSetValues(db, setValue); err != nil {
		http.Error(w, "Error updating allowance set values", http.StatusInternalServerError)
		return
	}

	responseData := struct {
		KReceipt float64 `json:"kReceipt"`
		Username string  `json:"username"`
	}{
		KReceipt: requestData.Amount,
		Username: username,
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

func UpdatePersonalHandler(w http.ResponseWriter, r *http.Request) {
	// Check basic authentication
	username, ok := basicAuth(w, r)
	if !ok {
		return
	}

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

	if err := Model.UpdatePersonalSetValues(db, setValue); err != nil {
		http.Error(w, "Error updating allowance set values", http.StatusInternalServerError)
		return
	}

	responseData := struct {
		PersonalDeduction float64 `json:"personalDeduction"`
		Username          string  `json:"username"`
	}{
		PersonalDeduction: requestData.Amount,
		Username:          username,
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
