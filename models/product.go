package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title		string             `bson:"title"`
	Description string			   `bson:"description"`
	Price 		string             `bson:"price"`
	Image    string             `bson:"image"`
}

