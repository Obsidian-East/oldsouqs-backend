package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItem struct {
	ProductID string `json:"productId" bson:"productId"`
	Quantity  int    `json:"quantity" bson:"quantity"`
}

type Cart struct {
	Id     primitive.ObjectID 	     `bson:"Id"`
	UserID string     `json:"userId" bson:"userId"`
	Items  []CartItem `json:"items" bson:"items"`
}
