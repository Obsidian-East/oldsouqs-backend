package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	FirstName   string             `bson:"first_name"`
	LastName    string             `bson:"last_name"`
	PhoneNumber string             `bson:"phone_number"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
}

