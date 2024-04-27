// File: Model/db.go

package Model

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// InitializeDB initializes the database connection using the provided URL
func InitializeDB() (*sql.DB, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		// If when testing, try loading from "../.env"
		if err := godotenv.Load("../.env"); err != nil {
			panic("Error loading .env file")
		}
	}

	// Retrieve the database URL from the environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	// Open a connection to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Check if the connection is successful
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	fmt.Println("Connected to the database")

	return db, nil
}
