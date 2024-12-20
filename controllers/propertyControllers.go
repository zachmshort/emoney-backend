package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zachmshort/monopoly-backend/config"
	"github.com/zachmshort/monopoly-backend/models"
	"go.mongodb.org/mongo-driver/bson"
)

func MortgageProperty(c *gin.Context) {}
func AddProperty(c *gin.Context)      {}
func RemoveProperty(c *gin.Context)   {}
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
