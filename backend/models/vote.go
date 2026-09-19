package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Vote struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	PollID    bson.ObjectID `bson:"pollId" json:"pollId"`
	Option    string        `bson:"option" json:"option"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
}
