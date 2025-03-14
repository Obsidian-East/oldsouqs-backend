package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"oldsouqs-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	var user models.User

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Println("Error decoding request body:", err)
		return
	}

	fmt.Println("Received user:", user)

	// Get users collection
	userCollection := db.Collection("users")

	// Check if user already exists
	var existingUser models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		fmt.Println("User already exists:", user.Email)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		fmt.Println("Error hashing password:", err)
		return
	}
	user.Password = string(hashedPassword)

	// Insert user
	result, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		fmt.Println("Error inserting user:", err)
		return
	}
	fmt.Println("Inserted user ID:", result.InsertedID)

	fmt.Println("User created successfully:", user.Email)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}
