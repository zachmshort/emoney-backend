package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Property struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	RoomID           primitive.ObjectID `bson:"roomId" json:"roomId"`
	PlayerID         primitive.ObjectID `bson:"playerId,omitempty" json:"playerId,omitempty"`
	PropertyIndex    int                `bson:"propertyIndex" json:"propertyIndex"`
	HouseCost        int                `bson:"houseCost" json:"houseCost"`
	Group            string             `bson:"group" json:"group"`
	Color            string             `bson:"color" json:"color"`
	Price            int                `bson:"price" json:"price"`
	Name             string             `bson:"name" json:"name"`
	IsMortgaged      bool               `bson:"isMortgaged" json:"isMortgaged"`
	DevelopmentLevel int                `bson:"developmentLevel" json:"developmentLevel"`
	RentPrices       []int              `bson:"rentPrices" json:"rentPrices"`
	Images           []string           `bson:"images" json:"images"`
}

type PropertyDetails []struct {
	PropertyID string
	Count      int
}

type PropertyPurchase struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	RoomID     primitive.ObjectID `bson:"roomId" json:"roomId"`
	BuyerID    primitive.ObjectID `bson:"buyerId" json:"buyerId"`
	PropertyID primitive.ObjectID `bson:"propertyId" json:"propertyId"`
	TimeStamp  time.Time          `bson:"timestamp" json:"timestamp"`
	Price      int                `bson:"price" json:"price"`
	Status     string             `bson:"status" json:"status"`
}
