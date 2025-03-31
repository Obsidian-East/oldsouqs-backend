package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty"`
	Sku			  string			   `bson:"sku"`
	Title		  string               `bson:"title"`
	TitleAr		  string               `bson:"titleAr"`
	Description   string			   `bson:"description"`
	DescriptionAr string			   `bson:"descriptionAr"`
	Price 		  float64              `bson:"price"`
	Image         string               `bson:"image"`
	Tag			  []string			   `bson:"tag"`
}

type Collection struct {
	ID 			      primitive.ObjectID     `bson:"_id,omitempty"`
	CollectionName    string				 `bson:"collectionName"`
	CollectionNameAr  string				 `bson:"collectionNameAr"`
	Description		  string				 `bson:"description"`
	DescriptionAr	  string				 `bson:"descriptionAr"`
	ProductIds		  []primitive.ObjectID   `bson:"productIds"`
	ShowCollection	  bool                   `bson:"showCollection"`
}

