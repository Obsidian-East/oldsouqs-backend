package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID        string             `bson:"orderId" json:"orderId"`
	PhoneNumber	   string             `bson:"phoneNumber" json:"phoneNumber"`
	UserID         string             `bson:"userId" json:"userId"`
	Location       string             `bson:"userLocation" json:"userLocation"`
	Items          []CartItem         `bson:"items" json:"items"`
	Subtotal       float64            `bson:"subtotal" json:"subtotal"`
	Total          float64            `bson:"total" json:"total"`
	Discounted     bool               `bson:"discounted" json:"discounted"`
	CreatedAt      time.Time          `bson:"creationDate" json:"creationDate"`
}
