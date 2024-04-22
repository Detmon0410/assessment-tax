package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the database URL from the environment
	dbURL := os.Getenv("DATABASE_URL")

	// Parse the database URL
	parts := strings.Split(dbURL, "://")
	if len(parts) != 2 {
		log.Fatal("Invalid DATABASE_URL format")
	}

	// Open a connection to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer db.Close()

	// Ping the database to check the connection status
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// Define the table name
	tableName := "allowance"

	// Query to select only the ID column from the specified table
	query := fmt.Sprintf("SELECT id FROM %s;", tableName)

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Error querying database: %v", err)
	}
	defer rows.Close()

	// Iterate over the result rows and log the IDs
	fmt.Printf("IDs from table '%s':\n", tableName)
	for rows.Next() {
		var recordID int
		err := rows.Scan(&recordID)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		fmt.Println(recordID)
	}

	// Check for errors during row iteration
	err = rows.Err()
	if err != nil {
		log.Fatalf("Error iterating over rows: %v", err)
	}

	// Log success message
	fmt.Println("Successfully logged all IDs from the table!")
}
