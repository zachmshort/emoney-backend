package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zachmshort/emoney-backend/config"
	"github.com/zachmshort/emoney-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateRoom(c *gin.Context) {

	var requestBody struct {
		PlayerName   string `json:"playerName" binding:"required"`
		RoomName     string `json:"roomName" binding:"required"`
		RoomCode     string `json:"roomCode" binding:"required"`
		PlayerColor  string `json:"playerColor" binding:"required"`
		StartingCash int    `json:"startingCash"`
		Houses       int    `json:"houses"`
		Hotels       int    `json:"hotels"`
	}

	var startingCash int
	if requestBody.StartingCash == 0 {
		startingCash = 1500
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := primitive.NewObjectID()
	playerID := primitive.NewObjectID()

	room := models.Room{
		ID:   roomID,
		Name: requestBody.RoomName,
		Code: requestBody.RoomCode,
		RoomRules: models.RoomRules{
			StartingCash: startingCash,
			MaxHouses:    requestBody.Houses,
			MaxHotels:    requestBody.Hotels,
		},
		FreeParking: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	banker := models.Player{
		ID:       playerID,
		RoomID:   roomID,
		IsBanker: true,
		IsActive: true,
		Balance:  requestBody.StartingCash,
		Name:     requestBody.PlayerName,
		Color:    requestBody.PlayerColor,
	}

	eventHistory := models.EventHistory{
		ID:        primitive.NewObjectID(),
		RoomID:    roomID,
		TimeStamp: time.Now(),
		Event:     fmt.Sprintf("%s created a new room, %s", banker.Name, requestBody.RoomName),
	}

	properties := make([]models.Property, len(config.DefaultProperties))
	var interfaceSlice []any
	for i := range config.DefaultProperties {
		properties[i] = models.Property{
			ID:               primitive.NewObjectID(),
			RoomID:           roomID,
			PropertyIndex:    i,
			Color:            config.DefaultProperties[i].Color,
			Price:            config.DefaultProperties[i].Price,
			Group:            config.DefaultProperties[i].Group,
			DevelopmentLevel: 0,
			HouseCost:        config.DefaultProperties[i].HouseCost,
			RentPrices:       config.DefaultProperties[i].RentPrices,
			IsMortgaged:      false,
			Name:             config.DefaultProperties[i].Name,
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

		eventHistoryColl := config.DB.Collection("EventHistory")
		_, err = eventHistoryColl.InsertOne(sc, eventHistory)
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
		"roomCode": room.Code,
		"playerId": banker.ID,
	})
}

func JoinRoom(c *gin.Context) {
	var requestBody struct {
		RoomCode    string `json:"roomCode" binding:"required"`
		PlayerName  string `json:"playerName" binding:"required"`
		PlayerColor string `json:"playerColor" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomColl := config.DB.Collection("Room")
	var room models.Room
	err := roomColl.FindOne(c, bson.M{"code": requestBody.RoomCode}).Decode(&room)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	playerColl := config.DB.Collection("Player")
	existingPlayerCount, _ := playerColl.CountDocuments(c, bson.M{
		"roomId": room.ID,
		"$or": []bson.M{
			{"name": requestBody.PlayerName},
			{"color": requestBody.PlayerColor},
		},
	})

	if existingPlayerCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Name or color already taken"})
		return
	}

	newPlayerID := primitive.NewObjectID()
	newPlayer := models.Player{
		ID:       newPlayerID,
		RoomID:   room.ID,
		IsBanker: false,
		IsActive: true,
		Balance:  room.RoomRules.StartingCash,
		Name:     requestBody.PlayerName,
		Color:    requestBody.PlayerColor,
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
		"roomCode": room.Code,
	})
}

func CheckIfRoomCodeExists(c *gin.Context) {
	code := c.Param("code")

	roomCollection := config.DB.Collection("Room")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	count, err := roomCollection.CountDocuments(ctx, bson.M{"code": code})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error while fetching existing rooms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": count > 0})
}
