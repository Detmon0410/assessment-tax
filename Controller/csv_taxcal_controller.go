package Controller

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
)

const personalDeduction = 60000

type TaxRecord struct {
	TotalIncome float64 `json:"totalIncome"`
	WHT         float64 `json:"wht"`
	Donation    float64 `json:"donation"`
	Personal    float64 `json:"personal"`
	Tax         float64 `json:"tax"`
}

func UploadCSVHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the request
	file, header, err := r.FormFile("taxFile")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check if the file name is "taxes.csv"
	if header.Filename != "taxes.csv" {
		http.Error(w, "Invalid file name. File name must be 'taxes.csv'", http.StatusBadRequest)
		return
	}

	// Read the file content
	fileBytes, err := csv.NewReader(file).ReadAll()
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
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
			WHT:         wht,
			Donation:    donation,
			Personal:    personalDeduction,
			Tax:         tax,
		}
		taxRecords = append(taxRecords, record)
	}

	// Convert tax records to JSON
	var taxes []map[string]interface{}
	for _, record := range taxRecords {
		tax := map[string]interface{}{
			"totalIncome": record.TotalIncome,
			"tax":         record.Tax,
		}
		taxes = append(taxes, tax)
	}
	responseData := map[string][]map[string]interface{}{"taxes": taxes}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Unable to convert to JSON", http.StatusInternalServerError)
		return
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
