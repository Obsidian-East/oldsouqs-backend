package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"oldsouqs-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func createUser(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	res, err := db.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	user.ID = res.InsertedID.(primitive.ObjectID)
	respondJSON(w, user)
}

func getUsers(w http.ResponseWriter, _ *http.Request, db *mongo.Database) {
	cursor, err := db.Collection("users").Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var users []models.User
	if err = cursor.All(context.TODO(), &users); err != nil {
		http.Error(w, "Error decoding users", http.StatusInternalServerError)
		return
	}

	respondJSON(w, users)
}

func getUserByID(w http.ResponseWriter, r *http.Request, db *mongo.Database, id string) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	err = db.Collection("users").FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	respondJSON(w, user)
}

func updateUser(w http.ResponseWriter, r *http.Request, db *mongo.Database, id string) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updateData bson.M
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	res, err := db.Collection("users").UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": updateData},
	)
	if err != nil || res.MatchedCount == 0 {
		http.Error(w, "User not found or update failed", http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"message": "User updated successfully"})
}

func deleteUser(w http.ResponseWriter, _ *http.Request, db *mongo.Database, id string) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	res, err := db.Collection("users").DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil || res.DeletedCount == 0 {
		http.Error(w, "User not found or deletion failed", http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"message": "User deleted successfully"})
}

// Utility to write JSON response
func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
