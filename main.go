package main

import (
	"file-sharing-system/models"
	"file-sharing-system/routes"
	"log"
)

func main() {
	// Load environment variables
	models.InitDB()

	r := routes.SetupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
