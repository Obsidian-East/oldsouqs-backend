package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"oldsouqs-backend/config"
	"oldsouqs-backend/routes"
)

func main() {
	// Get PORT from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if PORT is not set
	}

	// Example route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server is running on port %s!", port)
	})

	// Start the server
	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
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
