package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global DB variable
var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
	// Retrieve environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Build the Data Source Name (DSN) for PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Initialize the GORM database connection
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Optional: Automatically migrate your schema, if needed
	// err = DB.AutoMigrate(&YourModel{})
	// if err != nil {
	// 	log.Fatalf("failed to migrate database schema: %v", err)
	// }
}

// GetDB returns the global DB instance
func GetDB() *gorm.DB {
	return DB
}
