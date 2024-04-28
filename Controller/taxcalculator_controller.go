package Controller

import (
	"net/http"

	"github.com/Detmon0410/assessment-tax/Model"
	"github.com/labstack/echo/v4"
)

type Err struct {
	Message string `json:"message"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type TaxInput struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type TaxResponse struct {
	Tax       float64    `json:"tax"`
	TaxLevels []TaxLevel `json:"taxLevels"`
	TaxRefund float64    `json:"taxRefund,omitempty"`
}

// //////////// For Story: EXP04 ///////////
func CalculateTax(c echo.Context) error {

	db, err := Model.InitializeDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Error initializing database"})
	}
	defer db.Close()

	allowancesDB, err := Model.GetAllAllowances(db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Error fetching allowances from database"})
	}

	var personal, donationMax, kReceiptMax, kReceiptMin float64
	for _, allowance := range allowancesDB {
		switch allowance.AllowanceType {
		case "personal":
			personal = float64(allowance.SetValue)
		case "donation":
			donationMax = float64(allowance.SetValue)
		case "k-receipt":
			kReceiptMax = float64(allowance.SetValue)
			kReceiptMin = float64(allowance.Min)
		}
	}

	var input TaxInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	// Validate TotalIncome and WHT
	if input.TotalIncome <= 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "TotalIncome must be greater than 0"})
	}
	if input.WHT < 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "WHT cannot be less than 0"})
	}
	if input.WHT > input.TotalIncome {
		return c.JSON(http.StatusBadRequest, Err{Message: "WHT cannot be greater than TotalIncome"})
	}
	for _, allowance := range input.Allowances {
		switch allowance.AllowanceType {
		case "donation", "k-receipt":

		default:
			return c.JSON(http.StatusBadRequest, Err{Message: "invalid allowance type"})
		}
	}

	totalAllowance := 0.0
	for _, allowance := range input.Allowances {
		switch allowance.AllowanceType {
		case "donation":
			if allowance.Amount > donationMax {
				totalAllowance += donationMax
			} else {
				totalAllowance += allowance.Amount
			}
		case "k-receipt":
			if allowance.Amount < kReceiptMin {
				return c.JSON(http.StatusBadRequest, Err{Message: "k-receipt amount must be more than k-receipt min value"})
			}
			if allowance.Amount > kReceiptMax {
				allowance.Amount = kReceiptMax
			}
			totalAllowance += allowance.Amount
		}
	}

	taxableIncome := input.TotalIncome - totalAllowance - personal

	var tax float64
	var taxLevels []TaxLevel

	switch {
	case taxableIncome <= 150000:
		tax = 0
		taxLevels = append(taxLevels, TaxLevel{"0-150,000", 0})
		taxLevels = append(taxLevels, TaxLevel{"150,001-500,000", 0})
		taxLevels = append(taxLevels, TaxLevel{"500,001-1,000,000", 0})
		taxLevels = append(taxLevels, TaxLevel{"1,000,001-2,000,000", 0})
		taxLevels = append(taxLevels, TaxLevel{"2,000,001 ขึ้นไป", 0})
	case taxableIncome <= 500000:
		tax = (taxableIncome-150000)*0.1 - input.WHT
		taxLevels = append(taxLevels, TaxLevel{"0-150,000", 0})
		taxLevels = append(taxLevels, TaxLevel{"150,001-500,000", tax})
		taxLevels = append(taxLevels, TaxLevel{"500,001-1,000,000", 0})
		taxLevels = append(taxLevels, TaxLevel{"1,000,001-2,000,000", 0})
		taxLevels = append(taxLevels, TaxLevel{"2,000,001 ขึ้นไป", 0})
	case taxableIncome <= 1000000:
		tax = 35000 + (taxableIncome-500000)*0.15 - input.WHT
		taxLevels = append(taxLevels, TaxLevel{"0-150,000", 0})
		taxLevels = append(taxLevels, TaxLevel{"150,001-500,000", 0})
		taxLevels = append(taxLevels, TaxLevel{"500,001-1,000,000", tax})
		taxLevels = append(taxLevels, TaxLevel{"1,000,001-2,000,000", 0})
		taxLevels = append(taxLevels, TaxLevel{"2,000,001 ขึ้นไป", 0})
	case taxableIncome <= 2000000:
		tax = 135000 + (taxableIncome-1000000)*0.2 - input.WHT
		taxLevels = append(taxLevels, TaxLevel{"0-150,000", 0})
		taxLevels = append(taxLevels, TaxLevel{"150,001-500,000", 0})
		taxLevels = append(taxLevels, TaxLevel{"500,001-1,000,000", 0})
		taxLevels = append(taxLevels, TaxLevel{"1,000,001-2,000,000", tax})
		taxLevels = append(taxLevels, TaxLevel{"2,000,001 ขึ้นไป", 0})

	default:
		tax = 335000 + (taxableIncome-2000000)*0.35 - input.WHT
		taxLevels = append(taxLevels, TaxLevel{"0-150,000", 0})
		taxLevels = append(taxLevels, TaxLevel{"150,001-500,000", 0})
		taxLevels = append(taxLevels, TaxLevel{"500,001-1,000,000", 0})
		taxLevels = append(taxLevels, TaxLevel{"1,000,001-2,000,000", 0})
		taxLevels = append(taxLevels, TaxLevel{"2,000,001 ขึ้นไป", tax})
	}

	var taxRefund float64
	if tax < 0 {
		taxRefund = -tax
	}

	response := TaxResponse{Tax: tax, TaxLevels: taxLevels, TaxRefund: taxRefund}
	return c.JSON(http.StatusOK, response)
}
