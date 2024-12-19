package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomCode    primitive.ObjectID `bson:"roomCode" json:"roomCode"`
	BankerId    primitive.ObjectID `bson:"bankerId" json:"bankerId"`
	FreeParking int32              `bson:"freeParking" json:"freeParking"`
	IsActive    bool               `bson:"isActive" json:"isActive"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
