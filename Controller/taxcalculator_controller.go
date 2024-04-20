package Controller

import (
	"fmt"
	"net/http"

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
	Tax float64 `json:"tax"`
}

func CalculateTax(c echo.Context) error {
	input, err := getTaxInput()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	totalAllowance := 0.0
	for _, allowance := range input.Allowances {
		totalAllowance += allowance.Amount
	}
	taxableIncome := input.TotalIncome - totalAllowance - input.WHT
	var tax float64
	switch {
	case taxableIncome <= 150000:
		tax = 0
	case taxableIncome <= 500000:
		tax = (taxableIncome - 150000) * 0.1
	case taxableIncome <= 1000000:
		tax = 35000 + (taxableIncome-500000)*0.15
	case taxableIncome <= 2000000:
		tax = 135000 + (taxableIncome-1000000)*0.2
	default:
		tax = 335000 + (taxableIncome-2000000)*0.35
	}
	response := TaxResponse{Tax: tax}
	return c.JSON(http.StatusOK, response)
}

func getTaxInput() (TaxInput, error) {

	input := TaxInput{
		TotalIncome: 500000.0,
		WHT:         0.0,
		Allowances: []Allowance{
			{
				AllowanceType: "donation",
				Amount:        0,
			},
			{
				AllowanceType: "personal",
				Amount:        60000,
			},
		},
	}

	for _, allowance := range input.Allowances {
		switch allowance.AllowanceType {
		case "donation", "k-receipt", "personal":

		default:
			return TaxInput{}, fmt.Errorf("invalid allowance type")
		}
	}
	return input, nil
}
