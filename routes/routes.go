package routes

import (
	"net/http"

	"oldsouqs-backend/controllers"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(db *mongo.Database) *mux.Router {
	router := mux.NewRouter()

	// Auth routes
	router.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		controllers.SignupHandler(w, r, db)
	}).Methods("POST")

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

	// Get products by IDs
	router.HandleFunc("/products/ids", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetProductsByIDs(w, r, db)
	}).Methods("POST")

	// Collection routes
	router.HandleFunc("/collections", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controllers.GetCollections(w, r, db)
		} else if r.Method == http.MethodPost {
			controllers.CreateCollection(w, r, db)
		}
	}).Methods("GET", "POST")

	router.HandleFunc("/collections/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controllers.GetCollectionByID(w, r, db)
		} else if r.Method == http.MethodPut {
			controllers.UpdateCollection(w, r, db)
		} else if r.Method == http.MethodDelete {
			controllers.DeleteCollection(w, r, db)
		}
	}).Methods("GET", "PUT", "DELETE")

	// Get products by collection ID
	router.HandleFunc("/collections/{id}/products", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetProductsByCollection(w, r, db)
	}).Methods("GET")

	// Cart routes (RESTful)
	router.HandleFunc("/cart", func(w http.ResponseWriter, r *http.Request) {
		controllers.AddToCart(w, r, db)
	}).Methods(http.MethodPost)
	
	router.HandleFunc("/cart", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetCart(w, r, db)
	}).Methods(http.MethodGet)
	
	router.HandleFunc("/cart/{productId}", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateCartItem(w, r, db)
	}).Methods(http.MethodPut)
	
	router.HandleFunc("/cart/{productId}", func(w http.ResponseWriter, r *http.Request) {
		controllers.RemoveFromCart(w, r, db)
	}).Methods(http.MethodDelete)

	router.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllOrders(w, r, db)
	}).Methods("GET")
	
	router.HandleFunc("/orders/{orderId}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetOrder(w, r, db)
	}).Methods("GET")
	
	router.HandleFunc("/orders/{userId}", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateOrder(w, r, db)
	}).Methods("POST")
	
	router.HandleFunc("/orders/{orderId}", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateOrder(w, r, db)
	}).Methods("PUT")
	
	router.HandleFunc("/orders/{orderId}", func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteOrder(w, r, db)
	}).Methods("DELETE")
	
	
	return router
}
