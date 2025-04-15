package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"oldsouqs-backend/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateOrder(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	userID := mux.Vars(r)["userId"]

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	order.ID = primitive.NewObjectID()
	order.OrderID = fmt.Sprintf("OS%d", time.Now().Unix())
	order.CreatedAt = time.Now()
	order.UserID = userID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := db.Collection("orders").InsertOne(ctx, order); err != nil {
		http.Error(w, "Error creating order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func GetOrder(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	orderID := mux.Vars(r)["orderId"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var order models.Order
	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	found := db.Collection("orders").FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	if found != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func GetAllOrders(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection("orders").Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		http.Error(w, "Error parsing orders", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	orderID := mux.Vars(r)["orderId"]

	var updatedOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&updatedOrder); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"items":       updatedOrder.Items,
			"total":       updatedOrder.Total,
			"subtotal":    updatedOrder.Subtotal,
			"userId":      updatedOrder.UserID,
			"createdAt":   updatedOrder.CreatedAt,
			"location":    updatedOrder.Location,
			"hasDiscount": updatedOrder.Discounted,
		},
	}

	_, err = db.Collection("orders").UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		http.Error(w, "Error updating order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"message": "Order updated successfully"})
}

func DeleteOrder(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	orderID := mux.Vars(r)["orderId"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	_, err = db.Collection("orders").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		http.Error(w, "Error deleting order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"message": "Order deleted successfully"})
}
