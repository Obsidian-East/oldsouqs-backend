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
	
	// Product routes
	router.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controllers.GetProducts(w, r, db)
		} else if r.Method == http.MethodPost {
			controllers.CreateProduct(w, r, db)
		}
	}).Methods("GET", "POST")

	router.HandleFunc("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controllers.GetProduct(w, r, db)
		} else if r.Method == http.MethodPut {
			controllers.UpdateProduct(w, r, db)
		} else if r.Method == http.MethodDelete {
			controllers.DeleteProduct(w, r, db)
		}
	}).Methods("GET", "PUT", "DELETE")
	
	return router
}
