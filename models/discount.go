package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Discount struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type      string             `bson:"type" json:"type"` // "product" or "category"
	TargetID  primitive.ObjectID `bson:"targetId" json:"targetId"`
	Value     float64            `bson:"value" json:"value"` // % discount (0–100)
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
