package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventHistory struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	RoomID       primitive.ObjectID `bson:"roomId" json:"roomId"`
	TimeStamp    time.Time          `bson:"timestamp" json:"timestamp"`
	Event        string             `bson:"event" json:"event"`                          // description of what happened
	Type         string             `bson:"type" json:"type"`                            // "TRANSFER", "PROPERTY_PURCHASE", "MORTGAGE"
	FromPlayerID primitive.ObjectID `bson:"fromPlayerId, omitempty" json:"fromPlayerId"` // The player who initiated the action
	ToPlayerID   primitive.ObjectID `bson:"toPlayerId, omitempty" json:"toPlayerId"`     // The player who recieved
	Amount       int                `bson:"amount, omitempty" json:"amount"`
}
