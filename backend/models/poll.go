package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Poll struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Question  string        `bson:"question" json:"question"`
	Options   []string      `bson:"options" json:"options"`
	CreatedBy string        `bson:"createdBy" json:"createdBy"`
	IsActive  bool          `bson:"isActive" json:"isActive"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
}
