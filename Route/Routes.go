// File: Route/routes.go

package Route

import (
	"github.com/Detmon0410/assessment-tax/Controller"
	"github.com/labstack/echo/v4"
)

func GetRoutes() *echo.Echo {
	e := echo.New()

	// Define routes
	e.POST("/tax/calculations", Controller.CalculateTax)

	// Route for updating allowance
	e.POST("/admin/deductions/k-receipt", func(c echo.Context) error {
		// Call the controller function with the response writer and request from the context
		Controller.UpdateKReceiptHandler(c.Response(), c.Request())
		return nil
	})

	// Route for retrieving all allowances
	e.GET("/admin/allowances", func(c echo.Context) error {
		// Call the controller function to get all allowances
		Controller.GetAllAllowancesHandler(c.Response(), c.Request())
		return nil
	})

	return e
}
