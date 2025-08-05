package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Discount struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TargetType string             `bson:"targetType" json:"targetType"` // "product" or "collection"
	TargetID   primitive.ObjectID `bson:"targetId" json:"targetId"`     // productId or collectionId
	Percentage float64            `bson:"percentage" json:"percentage"` // 0–100
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}
