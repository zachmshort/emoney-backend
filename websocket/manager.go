package websocket

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/zachmshort/monopoly-backend/controllers"
	"github.com/zachmshort/monopoly-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomManager struct {
	clients map[string]map[*Client]bool
	mu      sync.RWMutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		clients: make(map[string]map[*Client]bool),
	}
}

func (rm *RoomManager) AddClient(client *Client) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.clients[client.Room] == nil {
		rm.clients[client.Room] = make(map[*Client]bool)
	}
	rm.clients[client.Room][client] = true
}

func (rm *RoomManager) RemoveClient(client *Client) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if _, ok := rm.clients[client.Room]; ok {
		delete(rm.clients[client.Room], client)
		if len(rm.clients[client.Room]) == 0 {
			delete(rm.clients, client.Room)
		}
	}
}

func (rm *RoomManager) Broadcast(room string, message Message) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	if clients, ok := rm.clients[room]; ok {
		for client := range clients {
			err := client.Conn.WriteJSON(message)
			if err != nil {
				client.Conn.Close()
				delete(clients, client)
			}
		}
	}
}

func (rm *RoomManager) handleTransfer(client *Client, message Message) error {
	log.Printf("Starting transfer handling for room: %s", client.Room)

	payload, ok := message.Payload.(map[string]interface{})
	log.Printf("Transfer payload received: %+v", payload)
	if !ok {
		return errors.New("invalid payload format")
	}

	amount, err := strconv.Atoi(payload["amount"].(string))
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	roomIdStr := payload["roomId"].(string)
	roomObjID, err := primitive.ObjectIDFromHex(roomIdStr)
	if err != nil {
		return fmt.Errorf("invalid room ID: %v", err)
	}

	transfer := models.Transfer{
		ID:        primitive.NewObjectID(),
		RoomID:    roomObjID,
		Amount:    amount,
		Reason:    payload["reason"].(string),
		Type:      payload["type"].(string),
		TimeStamp: time.Now(),
		Status:    models.TransferPending,
	}

	var transferErr error
	switch transfer.Type {
	case "SEND":
		fromID, err := primitive.ObjectIDFromHex(payload["fromPlayerId"].(string))
		if err != nil {
			return fmt.Errorf("invalid fromPlayerId: %w", err)
		}
		toID, err := primitive.ObjectIDFromHex(payload["toPlayerId"].(string))
		if err != nil {
			return fmt.Errorf("invalid toPlayerId: %w", err)
		}
		transfer.FromPlayerID = fromID
		transfer.ToPlayerID = toID
		transferErr = controllers.PlayerTransfer(transfer)
	case "REQUEST":
		transferErr = errors.New("request transfers not implemented yet")
	case "ADD", "SUBTRACT":
		transferErr = errors.New("bank transfers not implemented yet")
	default:
		transferErr = fmt.Errorf("invalid transfer type: %s", transfer.Type)
	}

	if transferErr != nil {
		return transferErr
	}

	transfer.Status = models.TransferCompleted
	log.Printf("Transfer successful, broadcasting update to room: %s", roomIdStr)
	var fromPlayer, toPlayer *models.Player

	fromPlayer, err = controllers.GetPlayer(transfer.FromPlayerID)
	if err != nil {
		log.Printf("Failed to get from player details: %v", err)
		return err
	}

	toPlayer, err = controllers.GetPlayer(transfer.ToPlayerID)
	if err != nil {
		log.Printf("Failed to get to player details: %v", err)
		return err
	}

	rm.Broadcast(client.Room, Message{
		Type: "TRANSFER",
		Payload: map[string]interface{}{
			"type":         "TRANSFER",
			"transfer":     transfer,
			"notification": fmt.Sprintf("%s just sent %s to %s ", fromPlayer.Name, strconv.Itoa(amount), toPlayer.Name),
		},
	})
	// fmt.Printf("%s just sent $%s to %s for %s", fromPlayer.Name, strconv.Itoa(amount), toPlayer.Name, transfer.Reason)
	fmt.Printf("BY ITSELF %s BY ITSELF", transfer.Reason)

	return nil

}

func (rm *RoomManager) handlePropertyPurchase(client *Client, message Message) error {
	log.Printf("Starting property purchase handling for room: %s", client.Room)

	payload, ok := message.Payload.(map[string]interface{})
	log.Printf("Property purchase payload received: %+v", payload)
	if !ok {
		return errors.New("invalid payload format")
	}

	priceFloat, ok := payload["price"].(float64)
	if !ok {
		return fmt.Errorf("invalid price format")
	}
	price := int(priceFloat)
	log.Printf("Processed price: %d", price)

	buyerID, err := primitive.ObjectIDFromHex(payload["buyerId"].(string))
	if err != nil {
		log.Printf("Invalid buyerId error: %v", err)
		return fmt.Errorf("invalid buyerId: %w", err)
	}
	log.Printf("Processed buyerId: %s", buyerID.Hex())

	propertyID, err := primitive.ObjectIDFromHex(payload["propertyId"].(string))
	if err != nil {
		log.Printf("Invalid propertyId error: %v", err)
		return fmt.Errorf("invalid propertyId: %w", err)
	}
	log.Printf("Processed propertyId: %s", propertyID.Hex())
	property, buyer, err := controllers.GetPropertyAndBuyer(propertyID, buyerID)
	if err != nil {
		log.Printf("Failed to get property or buyer details: %v", err)
		return err
	}

	log.Printf("Attempting to update property %s with new owner %s", propertyID.Hex(), buyerID.Hex())
	purchaseErr := controllers.PurchaseProperty(propertyID, buyerID, price)
	if purchaseErr != nil {
		log.Printf("Property update failed: %v", purchaseErr)
		return purchaseErr
	}
	log.Printf("Property update successful")

	log.Printf("Broadcasting update to room: %s", client.Room)
	rm.Broadcast(client.Room, Message{
		Type: "PURCHASE_PROPERTY",
		Payload: map[string]interface{}{
			"type":         "PURCHASE_PROPERTY",
			"propertyId":   propertyID.Hex(),
			"buyerId":      buyerID.Hex(),
			"price":        price,
			"propertyName": property.Name,
			"buyerName":    buyer.Name,
			"notification": fmt.Sprintf("%s has just purchased %s from the Bank", buyer.Name, property.Name),
		},
	})
	log.Printf("Broadcast complete to room: %s", client.Room)

	return nil
}
