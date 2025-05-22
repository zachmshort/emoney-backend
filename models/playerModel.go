package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	RoomID     primitive.ObjectID `bson:"roomId" json:"roomId"`
	Properties []Property         `bson:"-" json:"properties,omitempty"`
	IsBanker   bool               `bson:"isBanker" json:"isBanker"`
	IsActive   bool               `bson:"isActive" json:"isActive"`
	Balance    int                `bson:"balance" json:"balance"`
	Name       string             `bson:"name" json:"name"`
	Color      string             `bson:"color" json:"color"`
}
