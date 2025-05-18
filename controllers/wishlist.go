package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"oldsouqs-backend/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Add to Wishlist
func AddToWishlist(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	var item models.WishlistItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := r.URL.Query().Get("userId")
	if userID == "" || item.ProductID == "" {
		http.Error(w, "Missing userId or productId", http.StatusBadRequest)
		return
	}

	item.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wishlistCollection := db.Collection("wishlists")

	filter := bson.M{"userId": userID}
	update := bson.M{
		"$push": bson.M{"wishlistItems": item},
		"$setOnInsert": bson.M{
			"userId":    userID,
			"createdAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := wishlistCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Item added successfully",
	})
}

// Get Wishlist
func GetWishlist(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "Missing userId", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wishlist models.Wishlist
	collection := db.Collection("wishlists")
	err := collection.FindOne(ctx, bson.M{"userId": userID}).Decode(&wishlist)
	if err == mongo.ErrNoDocuments {
		json.NewEncoder(w).Encode([]models.WishlistItem{})
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(wishlist.WishlistItems)
}

// Remove item from wishlist
func RemoveFromWishlist(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	vars := mux.Vars(r)
	itemID := vars["itemId"]
	userID := r.URL.Query().Get("userId")

	if userID == "" || itemID == "" {
		http.Error(w, "Missing userId or itemId", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		http.Error(w, "Invalid itemId", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wishlistCollection := db.Collection("wishlists")

	filter := bson.M{"userId": userID}
	update := bson.M{"$pull": bson.M{"wishlistItems": bson.M{"_id": objID}}}

	_, err = wishlistCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, "Failed to remove wishlist item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Wishlist item removed")
}
