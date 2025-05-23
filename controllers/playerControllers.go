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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPlayersInRoom(c *gin.Context) {
	code := c.Param("code")
	playerId := c.Query("playerId")

	roomCollection := config.DB.Collection("Room")
	playerCollection := config.DB.Collection("Player")
	propertyCollection := config.DB.Collection("Property")
	eventHistoryCollection := config.DB.Collection("EventHistory")

	var room models.Room
	err := roomCollection.FindOne(c, bson.M{"code": code}).Decode(&room)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	query := bson.M{"roomId": room.ID}

	var existingPlayer *models.Player
	if playerId != "" {
		playerObjectID, err := primitive.ObjectIDFromHex(playerId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
			return
		}

		var player models.Player
		err = playerCollection.FindOne(c, bson.M{
			"_id":    playerObjectID,
			"roomId": room.ID,
		}).Decode(&player)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Player not found in this room",
				"isValid": false,
			})
			return
		}

		existingPlayer = &player
	}

	var players []models.Player
	cursor, err := playerCollection.Find(c, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get players"})
		return
	}
	if err = cursor.All(c, &players); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode players"})
		return
	}

	for i, player := range players {
		var properties []models.Property
		propCursor, err := propertyCollection.Find(c, bson.M{
			"roomId":   room.ID,
			"playerId": player.ID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get properties"})
			return
		}
		if err := propCursor.All(c, &properties); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode properties"})
			return
		}

		players[i].Properties = properties
	}

	var eventHistory []models.EventHistory
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor, err = eventHistoryCollection.Find(c, query, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get event history"})
		return
	}
	if err = cursor.All(c, &eventHistory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode event history"})
		return
	}

	response := gin.H{
		"players":      players,
		"room":         room,
		"eventHistory": eventHistory,
	}

	if existingPlayer != nil {
		response["existingPlayer"] = map[string]any{
			"id":      existingPlayer.ID.Hex(),
			"name":    existingPlayer.Name,
			"color":   existingPlayer.Color,
			"isValid": true,
		}
	}

	c.JSON(http.StatusOK, response)
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

	response := map[string]any{
		"player":     player,
		"properties": properties,
	}

	c.JSON(http.StatusOK, response)
}

func GetPlayer(buyerID primitive.ObjectID) (*models.Player, error) {
	var buyer models.Player
	err := config.DB.Collection("Player").FindOne(context.Background(), bson.M{"_id": buyerID}).Decode(&buyer)
	if err != nil {
		return nil, fmt.Errorf("failed to find buyer: %w", err)
	}

	return &buyer, nil
}

func UpdateFreeParkingBalance(roomID primitive.ObjectID, amount int, isAdd bool, playerID *primitive.ObjectID) error {
	session, err := config.DB.Client().StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(ctx mongo.SessionContext) (any, error) {
		freeParkingChange := amount
		if !isAdd {
			freeParkingChange = -amount
		}

		_, err := config.DB.Collection("Room").UpdateOne(
			ctx,
			bson.M{"_id": roomID},
			bson.M{"$inc": bson.M{"freeParking": freeParkingChange}},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update free parking: %w", err)
		}

		if !isAdd && playerID != nil {
			_, err = config.DB.Collection("Player").UpdateOne(
				ctx,
				bson.M{"_id": *playerID},
				bson.M{"$inc": bson.M{"balance": amount}},
			)
			if err != nil {
				return nil, fmt.Errorf("failed to update player balance: %w", err)
			}
		}

		return nil, nil
	})

	return err
}

func UpdatePlayerBalanceByBanker(roomID primitive.ObjectID, playerID primitive.ObjectID, amount int, isAdd bool) error {
	session, err := config.DB.Client().StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(ctx mongo.SessionContext) (any, error) {
		balanceChange := amount
		if !isAdd {
			balanceChange = -amount
		}

		_, err = config.DB.Collection("Player").UpdateOne(
			ctx,
			bson.M{"_id": playerID},
			bson.M{"$inc": bson.M{"balance": balanceChange}},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update player balance: %w", err)
		}

		return nil, nil
	})

	return err
}

func DeleteRoom(c *gin.Context) {
	code := c.Param("code")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	doc, err := config.DB.Collection("Room").DeleteOne(ctx, bson.M{"code": code})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting room"})
	}
	fmt.Println(doc)

	c.JSON(http.StatusOK, gin.H{"message": "Room deleted sucessfully"})

}
