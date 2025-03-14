package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"oldsouqs-backend/models"
	"oldsouqs-backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func validate(user models.User) string {
	// Check minimum length
	if len(user.Password) < 10 {
		return "Password must contain at least 10 characters"
	}

	// Check for at least one uppercase letter, one number, and one special character
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range user.Password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case (char >= '!' && char <= '/') || (char >= ':' && char <= '@') || (char >= '[' && char <= '`') || (char >= '{' && char <= '~'):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return "Password must contain at least one uppercase letter"
	}
	if !hasDigit {
		return "Password must contain at least one number"
	}
	if !hasSpecial {
		return "Password must contain at least one special character"
	}

	// Email validation
	emailRegex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, user.Email)
	if !matched {
		return "Invalid email format"
	}

	// Phone number validation
	phoneRegex := `^(?:\+961|00961)[0-9]{8}$`
	matched, _ = regexp.MatchString(phoneRegex, user.PhoneNumber)
	if !matched {
		return "Phone number must start with +961 or 00961 followed by 8 digits"
	}

	return ""
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

	valError := validate(user)
	if valError != "" {
		http.Error(w, valError, http.StatusBadRequest)
		fmt.Println("Error in validation:", valError)
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

func LoginHandler(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	var userInput models.User
	// Decode the login data (email and password)
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Fetch user from the database by email
	var user models.User
	collection := db.Collection("users")
	err = collection.FindOne(r.Context(), bson.M{"email": userInput.Email}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	if err != nil {
		http.Error(w, "Wrong Password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token (optional, if using JWT for session management)
	token, err := utils.GenerateJWT(user)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Return success with the JWT token
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message": "Login successful",
		"token":   token,
	}
	json.NewEncoder(w).Encode(response)
}