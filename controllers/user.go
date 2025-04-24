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

func validateUser(user models.User) string {
	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.PhoneNumber == "" || user.Location == "" {
		return "All fields are required"
	}
	return ""
}

func GetUsers(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	cursor, err := db.Collection("users").Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var users []models.User
	if err = cursor.All(context.TODO(), &users); err != nil {
		http.Error(w, "Error decoding users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func GetUserByID(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	err = db.Collection("users").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if validationError := validateUser(user); validationError != "" {
		http.Error(w, validationError, http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"first_name":   user.FirstName,
			"last_name":    user.LastName,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"location":     user.Location,
		},
	}

	_, err = db.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	_, err = db.Collection("users").DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
