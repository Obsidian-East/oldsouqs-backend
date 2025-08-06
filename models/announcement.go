// models/announcement.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Announcement struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Message   string             `bson:"message" json:"message"`
    CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
    UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
