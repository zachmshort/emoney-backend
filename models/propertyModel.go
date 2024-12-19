package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Property struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	RoomID      primitive.ObjectID `bson:"roomId" json:"roomId"`
	PlayerID    primitive.ObjectID `bson:"playerId" json:"playerId"`
	Name        string             `bson:"name" json:"name"`
	Houses      int                `bson:"houses" json:"houses"`
	Hotel       int                `bson:"hotel" json:"hotel"`
	Images      []string           `bson:"images" json:"images"`
	IsMortgaged bool               `bson:"isMortgaged" json:"isMortgaged"`
}
