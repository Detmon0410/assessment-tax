package Controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateTax_CustomInput(t *testing.T) {
	// Define the input JSON data
	inputData := []byte(`{
		"totalIncome": 500000.0,
		"wht": 0.0,
		"allowances": [
			{
				"allowanceType": "k-receipt",
				"amount": 200000.0
			},
			{
				"allowanceType": "donation",
				"amount": 100000.0
			}
		]
	}`)

	// Create a request with the input JSON data
	req, err := http.NewRequest(http.MethodPost, "/calculate-tax", bytes.NewBuffer(inputData))
	assert.NoError(t, err)

	// Set the request header
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to record the response
	rec := httptest.NewRecorder()

	// Call the CalculateTax handler function
	handler := http.HandlerFunc(CalculateTaxHandler)
	handler.ServeHTTP(rec, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Decode the response body
	var response TaxResponse
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Assert that the calculated tax matches the expected value
	expectedTax := 14000.0
	assert.Equal(t, expectedTax, response.Tax)
}
