package Route

import (
	"github.com/Detmon0410/assessment-tax/Controller"
	"github.com/labstack/echo/v4"
)

func GetRoutes() *echo.Echo {
	e := echo.New()
	e.POST("/tax/calculations", Controller.CalculateTax)

	return e
}
