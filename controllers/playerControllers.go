package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zachmshort/monopoly-backend/config"
	"github.com/zachmshort/monopoly-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetPlayersInRoom(c *gin.Context) {
	roomCode := c.Param("roomCode")

	collection := config.DB.Collection("Player")
	var players []models.Player
	cursor, err := collection.Find(c, bson.M{"roomId": roomCode})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get players"})
		return
	}

	if err = cursor.All(c, &players); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode players"})
		return
	}

	c.JSON(http.StatusOK, players)
}

func GetPlayerByDevice(c *gin.Context) {
	deviceId := c.Param("deviceId")

	collection := config.DB.Collection("Player")
	var player models.Player
	err := collection.FindOne(c, bson.M{"deviceId": deviceId}).Decode(&player)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	c.JSON(http.StatusOK, player)
}

func GetPlayerDetails(c *gin.Context) {
	playerId, _ := primitive.ObjectIDFromHex(c.Param("playerId"))

	playerColl := config.DB.Collection("Player")
	var player models.Player
	err := playerColl.FindOne(c, bson.M{"_id": playerId}).Decode(&player)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	propColl := config.DB.Collection("Property")
	var properties []models.Property
	cursor, err := propColl.Find(c, bson.M{"playerId": playerId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get properties"})
		return
	}

	if err = cursor.All(c, &properties); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode properties"})
		return
	}

	response := map[string]interface{}{
		"player":     player,
		"properties": properties,
	}

	c.JSON(http.StatusOK, response)
}
