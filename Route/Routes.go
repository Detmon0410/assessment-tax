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
	e.POST("/admin/deductions/personal", func(c echo.Context) error {
		// Call the controller function with the response writer and request from the context
		Controller.UpdatePersonalHandler(c.Response(), c.Request())
		return nil
	})

	// Route for uploading CSV and calculating tax
	e.POST("/tax/calculations/upload-csv", Controller.UploadCSVHandler)

	return e
}
