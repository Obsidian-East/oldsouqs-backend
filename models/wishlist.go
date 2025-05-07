package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type WishlistItem struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID string             `bson:"productId" json:"productId"`
}

type Wishlist struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        string             `bson:"userId" json:"userId"`
	CreatedAt     primitive.DateTime `bson:"createdAt,omitempty" json:"createdAt"`
	WishlistItems []WishlistItem     `bson:"wishlistItems" json:"wishlistItems"`
}
