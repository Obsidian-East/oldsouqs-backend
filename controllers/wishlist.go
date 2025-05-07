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
)

// Add item to wishlist
func AddToWishlist(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	var item models.WishlistItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "Missing userId", http.StatusBadRequest)
		return
	}

	// Assign a new ObjectID to the WishlistItem
	item.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wishlistCollection := db.Collection("wishlists")

	// Try to find an existing wishlist
	var existingWishlist models.Wishlist
	err := wishlistCollection.FindOne(ctx, bson.M{"userId": userID}).Decode(&existingWishlist)

	if err == mongo.ErrNoDocuments {
		// No wishlist yet, create a new one with the item
		newWishlist := models.Wishlist{
			ID:            primitive.NewObjectID(),
			UserID:        userID,
			CreatedAt:     primitive.NewDateTimeFromTime(time.Now()),
			WishlistItems: []models.WishlistItem{item},
		}

		_, insertErr := wishlistCollection.InsertOne(ctx, newWishlist)
		if insertErr != nil {
			http.Error(w, "Failed to create wishlist", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Wishlist created and item added",
		})
		return
	} else if err != nil {
		http.Error(w, "Error checking existing wishlist", http.StatusInternalServerError)
		return
	}

	// If wishlist exists, update it by pushing the item
	update := bson.M{
		"$push": bson.M{"wishlistItems": item},
	}

	_, updateErr := wishlistCollection.UpdateOne(ctx, bson.M{"userId": userID}, update)
	if updateErr != nil {
		http.Error(w, "Failed to add item to wishlist", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Item added to existing wishlist",
	})
}

// Get wishlist items for a specific user
func GetWishlist(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "Missing userId", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wishlistCollection := db.Collection("wishlists")

	var wishlist models.Wishlist
	err := wishlistCollection.FindOne(ctx, bson.M{"userId": userID}).Decode(&wishlist)
	if err == mongo.ErrNoDocuments {
		// No wishlist yet, return empty list instead of error
		json.NewEncoder(w).Encode([]models.WishlistItem{})
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
