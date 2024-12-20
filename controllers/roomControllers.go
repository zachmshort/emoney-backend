package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zachmshort/monopoly-backend/config"
	"github.com/zachmshort/monopoly-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateRoom(c *gin.Context) {

	var requestBody struct {
		Name  string `json:"name" binding:"required"`
		Code  string `json:"code" binding:"required"`
		Color string `json:"color" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
		IsBanker: true,
		IsActive: true,
		Balance:  1500,
		Name:     requestBody.Name,
		Color:    requestBody.Color,
	}

	properties := make([]models.Property, len(config.DefaultProperties))
	var interfaceSlice []interface{}
	for i := range config.DefaultProperties {
		properties[i] = models.Property{
			ID:            primitive.NewObjectID(),
			RoomID:        roomID,
			PropertyIndex: i,
			Houses:        0,
			Hotel:         0,
			HouseCost:     properties[i].HouseCost,
			IsMortgaged:   false,
			Images:        properties[i].Images,
		}
		interfaceSlice = append(interfaceSlice, properties[i])
	}

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

func JoinRoom(c *gin.Context) {
	var requestBody struct {
		RoomCode string `json:"roomCode" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Color    string `json:"color" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomColl := config.DB.Collection("Room")
	var room models.Room
	err := roomColl.FindOne(c, bson.M{"roomCode": requestBody.RoomCode}).Decode(&room)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	playerColl := config.DB.Collection("Player")
	existingPlayer, err := playerColl.CountDocuments(c, bson.M{
		"roomId": room.ID,
		"$or": []bson.M{
			{"name": requestBody.Name},
			{"color": requestBody.Color},
		},
	})

	if existingPlayer > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Name or color already taken"})
		return
	}

	newPlayerID := primitive.NewObjectID()
	newPlayer := models.Player{
		ID:       newPlayerID,
		RoomID:   room.ID,
		IsBanker: false,
		IsActive: true,
		Balance:  1500,
		Name:     requestBody.Name,
		Color:    requestBody.Color,
	}

	_, err = playerColl.InsertOne(c, newPlayer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join room"})
		return
	}

	var players []models.Player
	cursor, err := playerColl.Find(c, bson.M{"roomId": room.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get players"})
		return
	}

	if err = cursor.All(c, &players); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode players"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Successfully joined room",
		"playerId": newPlayerID.Hex(),
		"players":  players,
		"room":     room,
	})
}
