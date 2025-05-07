package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName   string             `bson:"first_name" json:"first_name"`
	LastName    string             `bson:"last_name" json:"last_name"`
	PhoneNumber string             `bson:"phonenumber" json:"phonenumber"`
	Location    string             `bson:"location" json:"location"`
	Email       string             `bson:"email" json:"email"`
	Password    string             `bson:"password" json:"-"`
  }
  

