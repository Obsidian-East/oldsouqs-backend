package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"oldsouqs-backend/config"
	"oldsouqs-backend/routes"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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

	// Wrap router with CORS middleware
	corsRouter := enableCORS(router)

	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, corsRouter))
}
