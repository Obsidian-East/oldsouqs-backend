package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Sku           string             `bson:"sku" json:"sku"`
	Title         string             `bson:"title" json:"title"`
	TitleAr       string             `bson:"titleAr" json:"titleAr"`
	Description   string             `bson:"description" json:"description"`
	DescriptionAr string             `bson:"descriptionAr" json:"descriptionAr"`
	Price         float64            `bson:"price" json:"price"`
	OriginalPrice *float64           `bson:"originalPrice,omitempty" json:"originalPrice,omitempty"` // Added for discount reversion
	Image         string             `bson:"image" json:"image"`
	Tag           []string           `bson:"tag" json:"tag"`
	Stock         int32              `bson:"stock" json:"stock"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Collection struct {
	ID               primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	CollectionName   string               `bson:"collectionName" json:"collectionName"`
	CollectionNameAr string               `bson:"collectionNameAr" json:"collectionNameAr"`
	Description      string               `bson:"description" json:"description"`
	DescriptionAr    string               `bson:"descriptionAr" json:"descriptionAr"`
	ProductIds       []primitive.ObjectID `bson:"productIds" json:"productIds"`
	ShowCollection   bool                 `bson:"showCollection" json:"showCollection"`
}
