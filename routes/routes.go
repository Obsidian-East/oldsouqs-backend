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

	// User routes
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetUsers(w, r, db)
	}).Methods("GET")

	router.HandleFunc("/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetUserByID(w, r, db)
	}).Methods("GET")

	router.HandleFunc("/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateUser(w, r, db)
	}).Methods("PUT")

	router.HandleFunc("/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteUser(w, r, db)
	}).Methods("DELETE")

	// Product routes
	router.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetProducts(w, r, db)
		case http.MethodPost:
			controllers.CreateProduct(w, r, db)
		}
	}).Methods("GET", "POST")

	router.HandleFunc("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetProduct(w, r, db)
		case http.MethodPut:
			controllers.UpdateProduct(w, r, db)
		case http.MethodDelete:
			controllers.DeleteProduct(w, r, db)
		}
	}).Methods("GET", "PUT", "DELETE")

	// Get products by IDs
	router.HandleFunc("/products/ids", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetProductsByIDs(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateProduct(w, r, db)
	}).Methods("POST")

	// Add Arabic routes
	router.HandleFunc("/ar/products", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetProducts(w, r, db)
	}).Methods("GET")

	router.HandleFunc("/ar/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetProduct(w, r, db)
	}).Methods("GET")

	// Collection routes
	router.HandleFunc("/collections", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetCollections(w, r, db)
		case http.MethodPost:
			controllers.CreateCollection(w, r, db)
		}
	}).Methods("GET", "POST")

	router.HandleFunc("/collections/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetCollectionByID(w, r, db)
		case http.MethodPut:
			controllers.UpdateCollection(w, r, db)
		case http.MethodDelete:
			controllers.DeleteCollection(w, r, db)
		}
	}).Methods("GET", "PUT", "DELETE")

	// Get products by collection ID
	router.HandleFunc("/collections/{id}/products", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetProductsByCollection(w, r, db)
	}).Methods("GET")

	router.HandleFunc("/ar/collections/{id}/products", func(w http.ResponseWriter, r *http.Request) {
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

	// order
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

	// Wishlist
	router.HandleFunc("/wishlist", func(w http.ResponseWriter, r *http.Request) {
		controllers.AddToWishlist(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/wishlist", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetWishlist(w, r, db)
	}).Methods("GET")

	router.HandleFunc("/wishlist/{itemId}", func(w http.ResponseWriter, r *http.Request) {
		controllers.RemoveFromWishlist(w, r, db)
	}).Methods("DELETE")

	// Sirv Handle
	router.HandleFunc("/api/upload", controllers.UploadImageToSirv).Methods("POST")

	// Discount routes
	controller := controllers.NewDiscountController(db)

	router.HandleFunc("/discounts", controller.CreateDiscount).Methods("POST")
	router.HandleFunc("/discounts", controller.GetDiscounts).Methods("GET")
	router.HandleFunc("/discounts/{id}", controller.UpdateDiscount).Methods("PUT")
	router.HandleFunc("/discounts/{id}", controller.DeleteDiscount).Methods("DELETE")

	// Announcement routes
	announcementController := controllers.NewAnnouncementController(db)
	
	router.HandleFunc("/announcements", announcementController.CreateAnnouncement).Methods("POST")
	router.HandleFunc("/announcements", announcementController.GetAnnouncements).Methods("GET")
	router.HandleFunc("/announcements/{id}", announcementController.UpdateAnnouncement).Methods("PUT")
	router.HandleFunc("/announcements/{id}", announcementController.DeleteAnnouncement).Methods("DELETE")

	return router
}
