package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"oldsouqs-backend/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateDiscount adds a new discount to the database
func CreateDiscount(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	var discount models.Discount
	if err := json.NewDecoder(r.Body).Decode(&discount); err != nil {
		http.Error(w, "Invalid discount data", http.StatusBadRequest)
		return
	}

	discount.ID = primitive.NewObjectID()
	_, err := db.Collection("discounts").InsertOne(context.TODO(), discount)
	if err != nil {
		http.Error(w, "Failed to create discount", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(discount)
}

// GetDiscounts returns all discounts
func GetDiscounts(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	cursor, err := db.Collection("discounts").Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch discounts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var discounts []models.Discount
	if err := cursor.All(context.TODO(), &discounts); err != nil {
		http.Error(w, "Failed to parse discounts", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(discounts)
}

// UpdateDiscount modifies an existing discount
func UpdateDiscount(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid discount ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	delete(updates, "_id")

	result, err := db.Collection("discounts").UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": updates},
	)
	if err != nil || result.MatchedCount == 0 {
		http.Error(w, "Failed to update discount", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(bson.M{"updated": true})
}

// DeleteDiscount removes a discount by ID
func DeleteDiscount(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid discount ID", http.StatusBadRequest)
		return
	}

	result, err := db.Collection("discounts").DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil || result.DeletedCount == 0 {
		http.Error(w, "Failed to delete discount", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(bson.M{"deleted": true})
}
