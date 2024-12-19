package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zachmshort/monopoly-backend/config"
	"github.com/zachmshort/monopoly-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateRoom(c *gin.Context) {

	var requestBody struct {
		Name     string `json:"name" binding:"required"`
		DeviceID string `json:"deviceId" binding:"required"`
		Code     string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Change deviceID to be just a string
	roomID := primitive.NewObjectID()
	playerID := primitive.NewObjectID()

	room := models.Room{
		ID:          roomID,
		RoomCode:    requestBody.Code,
		BankerId:    playerID,
		FreeParking: 0,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	banker := models.Player{
		ID:       playerID,
		RoomID:   roomID,
		DeviceID: requestBody.DeviceID, // Just use the string directly
		IsBanker: true,
		IsActive: true,
		Balance:  1500,
		Name:     requestBody.Name,
		Color:    "#FF0000",
	}
	// Initialize properties
	properties := make([]models.Property, len(config.DefaultProperties))
	var interfaceSlice []interface{}
	for i := range config.DefaultProperties {
		properties[i] = models.Property{
			ID:            primitive.NewObjectID(),
			RoomID:        roomID,
			PropertyIndex: i,
			Houses:        0,
			Hotel:         0,
			IsMortgaged:   false,
		}
		interfaceSlice = append(interfaceSlice, properties[i])
	}

	// Start MongoDB transaction
	session, err := config.DB.Client().StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start session"})
		return
	}
	defer session.EndSession(c)

	err = session.StartTransaction()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	err = mongo.WithSession(c, session, func(sc mongo.SessionContext) error {
		roomColl := config.DB.Collection("Room")
		_, err := roomColl.InsertOne(sc, room)
		if err != nil {
			return err
		}

		playerColl := config.DB.Collection("Player")
		_, err = playerColl.InsertOne(sc, banker)
		if err != nil {
			return err
		}

		propColl := config.DB.Collection("Property")
		_, err = propColl.InsertMany(sc, interfaceSlice)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})

	if err != nil {
		session.AbortTransaction(c)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roomId":   room.ID,
		"roomCode": room.RoomCode,
		"playerId": banker.ID,
	})
}
