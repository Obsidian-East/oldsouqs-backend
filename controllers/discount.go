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

var discountCollection *mongo.Collection

func InitDiscountController(db *mongo.Database) {
	discountCollection = db.Collection("discounts")
}

func CreateDiscount(w http.ResponseWriter, r *http.Request) {
	var discount models.Discount
	if err := json.NewDecoder(r.Body).Decode(&discount); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if discount.Percentage < 0 || discount.Percentage > 100 {
		http.Error(w, "Invalid percentage value", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := discountCollection.InsertOne(ctx, discount)
	if err != nil {
		http.Error(w, "Failed to create discount", http.StatusInternalServerError)
		return
	}
	discount.ID = res.InsertedID.(primitive.ObjectID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(discount)
}

func GetDiscounts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := discountCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch discounts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var discounts []models.Discount
	if err := cursor.All(ctx, &discounts); err != nil {
		http.Error(w, "Error reading discounts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(discounts)
}

func UpdateDiscount(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	discountID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updated models.Discount
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"targetType": updated.TargetType,
			"targetId":   updated.TargetID,
			"percentage": updated.Percentage,
		},
	}

	_, err = discountCollection.UpdateOne(ctx, bson.M{"_id": discountID}, update)
	if err != nil {
		http.Error(w, "Failed to update discount", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteDiscount(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	discountID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = discountCollection.DeleteOne(ctx, bson.M{"_id": discountID})
	if err != nil {
		http.Error(w, "Failed to delete discount", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
