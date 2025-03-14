package utils

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"oldsouqs-backend/models"
	"os"
	"log"
)

// Secret key for signing the JWT token (Should be moved to an environment variable)
var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// GenerateJWT generates a JWT token for a valid user
func GenerateJWT(user models.User) (string, error) {
	// Create a new token object with the signing method and claims
	claims := &jwt.StandardClaims{
		Subject:   user.ID.Hex(), // User ID as the subject
		Issuer:    "OldSouqsApp", // You can add an app name here
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	}

	// Create the token with the specified claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using the secret key
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println("Error signing the token:", err)
		return "", err
	}

	return signedToken, nil
}
