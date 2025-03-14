package routes

import (
	"net/http"

	"oldsouqs-backend/controllers"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(db *mongo.Database) *mux.Router {
	router := mux.NewRouter()
	
	// Signup route
	router.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		controllers.SignupHandler(w, r, db)
	}).Methods("POST")

	// Login route
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controllers.LoginHandler(w, r, db)
	}).Methods("POST")

	return router
}
