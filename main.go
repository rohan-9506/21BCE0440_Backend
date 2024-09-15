package main

import (
	"file-sharing-system/api"
	"file-sharing-system/models"
	"file-sharing-system/routes"
	"log"
)

func main() {
	// Load environment variables
	models.InitDB()
	api.InitRedis()

	// Run migrations to ensure the database schema is up-to-date
	if err := models.GetDB().AutoMigrate(&models.File{}); err != nil {
		log.Fatalf("could not run migrations: %v", err)
	}

	// Setup and start the router
	r := routes.SetupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
