package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"oldsouqs-backend/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 🛒 Add item to cart
func AddToCart(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	var item models.CartItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "Missing userId", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cartCollection := db.Collection("carts")

	filter := bson.M{"userId": userID}
	update := bson.M{"$push": bson.M{"items": item}}
	opts := options.Update().SetUpsert(true)

	_, err := cartCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		http.Error(w, "Failed to add to cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Item added to cart")
}

// 📦 Get cart items for a specific user
func GetCart(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "Missing userId", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cartCollection := db.Collection("carts")

	var cart models.Cart
	err := cartCollection.FindOne(ctx, bson.M{"userId": userID}).Decode(&cart)
	if err != nil {
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(cart.Items)
}

// 🔄 Update item quantity
func UpdateCartItem(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	vars := mux.Vars(r)
	productID := vars["productId"]
	userID := r.URL.Query().Get("userId")

	if userID == "" || productID == "" {
		http.Error(w, "Missing userId or productId", http.StatusBadRequest)
		return
	}

	var body struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cartCollection := db.Collection("carts")

	filter := bson.M{"userId": userID, "items.productId": productID}
	update := bson.M{"$set": bson.M{"items.$.quantity": body.Quantity}}

	_, err := cartCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, "Failed to update cart item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Cart item updated")
}

// ❌ Remove item from cart
func RemoveFromCart(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	vars := mux.Vars(r)
	productID := vars["productId"]
	userID := r.URL.Query().Get("userId")

	if userID == "" || productID == "" {
		http.Error(w, "Missing userId or productId", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cartCollection := db.Collection("carts")

	filter := bson.M{"userId": userID}
	update := bson.M{"$pull": bson.M{"items": bson.M{"productId": productID}}}

	_, err := cartCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, "Failed to remove cart item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Cart item removed")
}
