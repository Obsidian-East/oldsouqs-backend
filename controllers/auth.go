package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"oldsouqs-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func validate(user models.User) error {
	// Password validation (at least 10 characters, 1 special, 1 uppercase, 1 number)
	passwordRegex := `^(?=.*[A-Z])(?=.*[0-9])(?=.*[\W_]).{10,}$`
	matched, _ := regexp.MatchString(passwordRegex, user.Password)
	if !matched {
		return errors.New("password must contain at least 10 characters, 1 special character, 1 uppercase letter, and 1 number")
	}

	// Email validation (basic format validation)
	emailRegex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	matched, _ = regexp.MatchString(emailRegex, user.Email)
	if !matched {
		return errors.New("invalid email format")
	}

	// Phone number validation (starting with +961 or 00961)
	phoneRegex := `^(?:\+961|00961)[0-9]{8,}$`
	matched, _ = regexp.MatchString(phoneRegex, user.PhoneNumber)
	if !matched {
		return errors.New("phone number must start with +961 or 00961 followed by 8 digits")
	}

	return nil
}

func SignupHandler(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	var user models.User

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Println("Error decoding request body:", err)
		return
	}

	err = validate(user)
	if err != nil {
		http.Error(w, "Error Validating your request", http.StatusBadRequest)
		fmt.Println("Error in validation:", err)
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
