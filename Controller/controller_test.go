package Controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Detmon0410/assessment-tax/Model"
	"github.com/labstack/echo/v4"
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

// //////////////////// Test Admin Authentication //////////////////////
func TestBasicAuth(t *testing.T) {
	// Mock a request with basic authentication
	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth("adminTax", "admin!")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the basicAuth function and assert the result
	username, ok := basicAuth(rr, req)
	assert.True(t, ok)
	assert.Equal(t, "adminTax", username)
}

func TestGetAllAllowancesHandler(t *testing.T) {
	// Create a request
	req := httptest.NewRequest("GET", "/get-all-allowances", nil)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	GetAllAllowancesHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body and perform assertions
	var allowances []Model.Allowance
	json.NewDecoder(rr.Body).Decode(&allowances)
	assert.NotEmpty(t, allowances)
}

func TestUploadCSVHandler(t *testing.T) {
	// Test cases
	tests := []struct {
		Name         string
		RequestBody  string
		ExpectedCode int
		ExpectedBody string
	}{
		{
			Name:         "ValidInput",
			RequestBody:  "totalIncome,wht,donation\n500000,0,0\n600000,40000,20000\n750000,50000,15000",
			ExpectedCode: http.StatusOK,
			ExpectedBody: `{"taxes":[{"totalIncome":500000,"tax":29000},{"totalIncome":600000,"tax":33000},{"totalIncome":750000,"tax":53750}]}`,
		},
		// Add more test cases for edge cases, invalid input, etc.
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// Create a new echo instance
			e := echo.New()

			// Create a buffer to write the form data
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			// Write the CSV data to the form field
			formField, err := writer.CreateFormFile("taxFile", "taxes.csv")
			if err != nil {
				t.Fatal(err)
			}
			csvData := strings.NewReader(test.RequestBody)
			io.Copy(formField, csvData)

			// Close the multipart writer
			writer.Close()

			// Create a request object with the multipart form data
			req := httptest.NewRequest(http.MethodPost, "/upload", body)
			req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())

			// Create a response recorder
			rec := httptest.NewRecorder()

			// Create a context from the request and recorder
			c := e.NewContext(req, rec)

			// Call the handler function
			if err := UploadCSVHandler(c); err != nil {
				t.Fatal(err)
			}

			// Check the response status code
			if rec.Code != test.ExpectedCode {
				t.Errorf("Expected status code %d but got %d", test.ExpectedCode, rec.Code)
			}

			// Check the response body
			actualBody := rec.Body.String()
			if actualBody != test.ExpectedBody {
				t.Errorf("Expected body %q but got %q", test.ExpectedBody, actualBody)
				// Print the response body for further inspection
				fmt.Println("Response body:", actualBody)
			}
		})
	}
}
