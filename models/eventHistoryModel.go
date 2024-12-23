package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventHistory struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	RoomID    primitive.ObjectID `bson:"roomId" json:"roomId"`
	TimeStamp time.Time          `bson:"timestamp" json:"timestamp"`
	Event     string             `bson:"event" json:"event"` // description of what happened
	EventType []string           `bson:"eventType"`
}
