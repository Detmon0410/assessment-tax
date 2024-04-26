package Controller

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

const personalDeduction = 60000

type TaxRecord struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
}

type TaxResponseCSV struct {
	Taxes []TaxRecord `json:"taxes"`
}

// UploadCSVHandler handles CSV file upload and tax calculation
func UploadCSVHandler(c echo.Context) error {
	// Parse the multipart form
	err := c.Request().ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		return c.String(http.StatusBadRequest, "Unable to parse form")
	}

	// Get the file from the request
	taxFile, taxHeader, err := c.Request().FormFile("taxFile")
	if err != nil {
		return c.String(http.StatusBadRequest, "Unable to get file")
	}
	defer taxFile.Close()

	// Check if the file name is "taxes.csv"
	if taxHeader.Filename != "taxes.csv" {
		return c.String(http.StatusBadRequest, "Invalid file name. File name must be 'taxes.csv'")
	}

	// Read the file content
	fileBytes, err := csv.NewReader(taxFile).ReadAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Unable to read file")
	}

	// Check if the file is empty
	if len(fileBytes) == 0 {
		return c.String(http.StatusBadRequest, "Empty file")
	}

	// Process CSV data
	var taxRecords []TaxRecord
	for i, row := range fileBytes {
		if i == 0 {
			// Skip header row
			continue
		}

		// Parse CSV values
		totalIncome, _ := strconv.ParseFloat(row[0], 64)
		wht, _ := strconv.ParseFloat(row[1], 64)
		donation, _ := strconv.ParseFloat(row[2], 64)

		// Calculate taxable income
		taxableIncome := totalIncome - donation - wht - personalDeduction
		if taxableIncome < 0 {
			taxableIncome = 0 // Ensure taxable income is non-negative
		}

		// Calculate tax
		var tax float64
		switch {
		case taxableIncome <= 150000:
			tax = 0
		case taxableIncome <= 500000:
			tax = (taxableIncome - 150000) * 10 / 100
		case taxableIncome <= 1000000:
			tax = 35000 + (taxableIncome-500000)*15/100
		case taxableIncome <= 2000000:
			tax = 125000 + (taxableIncome-1000000)*20/100
		default:
			tax = 325000 + (taxableIncome-2000000)*35/100
		}

		// Create tax record
		record := TaxRecord{
			TotalIncome: totalIncome,
			Tax:         tax,
		}
		taxRecords = append(taxRecords, record)
	}

	// Create response object
	response := TaxResponseCSV{
		Taxes: taxRecords,
	}

	// Convert response object to JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Unable to convert to JSON")
	}

	// Write JSON response
	return c.JSONBlob(http.StatusOK, jsonData)
}
