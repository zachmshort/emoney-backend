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

	rm.Broadcast(client.Room, Message{
		Type: "GAME_STATE_UPDATE",
		Payload: map[string]interface{}{
			"type":     "TRANSFER_COMPLETE",
			"transfer": transfer,
		},
	})
	log.Printf("Broadcast complete to room: %s", client.Room)

	return nil

}
