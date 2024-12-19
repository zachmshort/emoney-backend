package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Property struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	RoomID        primitive.ObjectID `bson:"roomId" json:"roomId"`
	PlayerID      primitive.ObjectID `bson:"playerId,omitempty" json:"playerId,omitempty"`
	PropertyIndex int                `bson:"propertyIndex" json:"propertyIndex"`
	Houses        int                `bson:"houses" json:"houses"`
	Hotel         int                `bson:"hotel" json:"hotel"`
	IsMortgaged   bool               `bson:"isMortgaged" json:"isMortgaged"`
	Images        []string           `bson:"images" json:"images"`
}
