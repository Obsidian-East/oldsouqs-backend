package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Sku			  string			   `bson:"sku"`
	Title		  string             `bson:"title"`
	TitleAr		  string             `bson:"titleAr"`
	Description   string			   `bson:"description"`
	DescriptionAr string			   `bson:"descriptionAr"`
	Price 		  float64            `bson:"price"`
	Image         string             `bson:"image"`
}

