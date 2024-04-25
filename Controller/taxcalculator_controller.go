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

type TaxResponse struct {
	Tax       float64 `json:"tax"`
	TaxRefund float64 `json:"tax_refund,omitempty"`
}

// //////////// For Story: EXP03 ///////////
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
		case "k_receipt":
			kReceiptMax = float64(allowance.SetValue)
			kReceiptMin = float64(allowance.Min)
		}
	}

	var input TaxInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
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
	switch {
	case taxableIncome <= 150000:
		tax = 0
	case taxableIncome <= 500000:
		tax = (taxableIncome-150000)*0.1 - input.WHT
	case taxableIncome <= 1000000:
		tax = 35000 + (taxableIncome-500000)*0.15 - input.WHT
	case taxableIncome <= 2000000:
		tax = 135000 + (taxableIncome-1000000)*0.2 - input.WHT
	default:
		tax = 335000 + (taxableIncome-2000000)*0.35 - input.WHT
	}

	if tax < 0 {
		taxRefund := -tax
		response := TaxResponse{TaxRefund: taxRefund}
		return c.JSON(http.StatusOK, response)
	}

	response := TaxResponse{Tax: tax}
	return c.JSON(http.StatusOK, response)
}
