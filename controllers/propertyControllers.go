package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zachmshort/monopoly-backend/config"
	"github.com/zachmshort/monopoly-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MortgageProperty(c *gin.Context) {}
func AddProperty(c *gin.Context)      {}
func RemoveProperty(c *gin.Context)   {}

func PurchaseProperty(propertyID, buyerID primitive.ObjectID, price int) error {
	_, err := config.DB.Collection("Property").UpdateOne(
		context.Background(),
		bson.M{"_id": propertyID},
		bson.M{"$set": bson.M{"playerId": buyerID}},
	)
	if err != nil {
		return err
	}

	_, err = config.DB.Collection("Player").UpdateOne(
		context.Background(),
		bson.M{"_id": buyerID},
		bson.M{"$inc": bson.M{"balance": -price}},
	)
	return err
}

func GetPropertyAndBuyer(propertyID, buyerID primitive.ObjectID) (*models.Property, *models.Player, error) {
	var property models.Property
	err := config.DB.Collection("Property").FindOne(context.Background(), bson.M{"_id": propertyID}).Decode(&property)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find property: %w", err)
	}

	var buyer models.Player
	err = config.DB.Collection("Player").FindOne(context.Background(), bson.M{"_id": buyerID}).Decode(&buyer)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find buyer: %w", err)
	}

	return &property, &buyer, nil
}

func AssignOwnerShipProperty(propertyID, buyerID primitive.ObjectID) error {
	_, err := config.DB.Collection("Property").UpdateOne(
		context.Background(),
		bson.M{"_id": propertyID},
		bson.M{"$set": bson.M{"playerId": buyerID}},
	)
	return err
}

func GetAvailableProperties(c *gin.Context) {
	roomCode := c.Param("roomCode")

	roomCollection := config.DB.Collection("Room")
	var room models.Room
	err := roomCollection.FindOne(c, bson.M{"roomCode": roomCode}).Decode(&room)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	propertyCollection := config.DB.Collection("Property")

	query := bson.M{
		"roomId":   room.ID,
		"playerId": nil,
	}

	var availableProperties []models.Property
	cursor, err := propertyCollection.Find(c, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve available properties"})
		return
	}

	if err = cursor.All(c, &availableProperties); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode available properties"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"availableProperties": availableProperties,
		"roomId":              room.ID,
	})
}
