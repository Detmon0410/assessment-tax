package main

import (
	"github.com/Detmon0410/assessment-tax/Route"

	"log"
)

func main() {
	// Database

	// Starting server
	echo := Route.GetRoutes()
	err := echo.Start((":8080"))
	if err != nil {
		log.Fatal(err)
	}

}
