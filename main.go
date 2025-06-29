package main

import (
	"log"
	"fmt"
	"net/http"
	"os"

	"oldsouqs-backend/config"
	"oldsouqs-backend/controllers" // Import your controllers package
	"oldsouqs-backend/routes"

	"github.com/joho/godotenv"
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
	fmt.Println("Sirv Uploader application started.")

	// You can test GetSirvToken directly:
	token, err := controllers.GetSirvToken() // Changed to GetSirvToken()
	if err != nil {
		log.Fatalf("Fatal error getting Sirv token: %v", err)
	}
	fmt.Printf("Successfully retrieved Sirv token: %s\n", token[:10]+"...\n")

	// For a complete runnable example, let's simulate an upload if a test file exists.
	dummyFilePath := "test_image.txt"
	err = os.WriteFile(dummyFilePath, []byte("This is a dummy image file content."), 0644)
	if err != nil {
		log.Printf("Could not create dummy file for testing: %v", err)
	} else {
		fmt.Printf("Dummy file '%s' created for testing upload.\n", dummyFilePath)
		// Simulate upload
		err = controllers.UploadToSirv(dummyFilePath, "my_test_document.txt", token) // Changed to UploadToSirv()
		if err != nil {
			log.Printf("Simulated upload failed: %v", err)
		} else {
			fmt.Println("Simulated upload successful!")
		}
		os.Remove(dummyFilePath) // Clean up dummy file
	}

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