package main

import (
	"fmt"
	"log"
	"net/http"

	"oldsouqs-backend/config"
	"oldsouqs-backend/routes"
)

func main() {
	// Initialize database connection
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Pass database instance to routes
	router := routes.SetupRoutes(db)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
