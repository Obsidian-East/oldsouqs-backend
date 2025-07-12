package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID 				primitive.ObjectID    `bson:"_id,omitempty" json:"-"`
	Sku 			string				  `bson:"sku" json:"sku"`
	Title 			string 				  `bson:"title" json:"title"`
	TitleAr 		string 			      `bson:"titleAr" json:"titleAr"`
	Description 	string                `bson:"description" json:"description"`
	DescriptionAr   string                `bson:"descriptionAr" json:"descriptionAr"`
	Price 			float64               `bson:"price" json:"price"`
	Image 			string                `bson:"image" json:"image"`
	Tag 			[]string              `bson:"tag" json:"tag"`
	Stock 			int32                 `bson:"stock" json:"stock"`
	CreatedAt 		time.Time             `bson:"createdAt" json:"createdAt"`
	UpdatedAt 		time.Time             `bson:"updatedAt" json:"updatedAt"`
}

type Collection struct {
	ID               primitive.ObjectID   `bson:"_id,omitempty"`
	CollectionName   string               `bson:"collectionName"`
	CollectionNameAr string               `bson:"collectionNameAr"`
	Description      string               `bson:"description"`
	DescriptionAr    string               `bson:"descriptionAr"`
	ProductIds       []primitive.ObjectID `bson:"productIds"`
	ShowCollection   bool                 `bson:"showCollection"`
}
