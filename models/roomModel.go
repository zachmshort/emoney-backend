package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	RoomCode    string             `bson:"roomCode" json:"roomCode"`
	FreeParking int                `bson:"freeParking" json:"freeParking"`
	RoomRules   RoomRules          `bson:"roomRules" json:"roomRules"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type RoomRules struct {
	StartingCash int `bson:"startingCash" json:"startingCash"`
	MaxHouses    int `bson:"maxHouses" json:"maxHouses"`
	MaxHotels    int `bson:"maxHotels" json:"maxHotels"`
}
