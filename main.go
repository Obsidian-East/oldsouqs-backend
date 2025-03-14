package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"oldsouqs-backend/config"
	"oldsouqs-backend/routes"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	
	// Get PORT from environment variable or default to 10000
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000" // Default port to match Render's assigned port
	}

	// Initialize database connection
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Pass database instance to routes
	router := routes.SetupRoutes(db)

	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
