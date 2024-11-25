package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/yourusername/go-user-service/api"
)

func main() {
	// Database connection
	var err error

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	api.Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer api.Db.Close()

	// Create users table if it doesn't exist
	createTable()

	// Set up router
	router := api.SetupRoutes()

	// Start server
	fmt.Println("Server is running on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func createTable() {
	query := `
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        email TEXT NOT NULL UNIQUE,
        name TEXT NOT NULL,
        password TEXT NOT NULL
    );
    `
	_, err := api.Db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	} else {
		log.Println("Table 'users' created or already exists.")
	}
}
