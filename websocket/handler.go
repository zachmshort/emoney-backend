package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		switch origin {
		case "http://localhost:3000",
			"https://emoney.club",
			"https://www.emoney.club":
			return true
		default:
			return false
		}
	},
}
var (
	Manager = NewRoomManager()
)

func HandleWebSocket(c *gin.Context) {
	roomCode := c.Param("code")
	log.Printf("WebSocket connection attempt for room: %s", roomCode)

	if roomCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room code required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		Conn:     conn,
		Room:     roomCode,
		PlayerID: "",
	}

	log.Printf("New client connected to room %s", roomCode)
	Manager.AddClient(client)

	defer func() {
		log.Printf("Client disconnecting from room %s", roomCode)
		Manager.RemoveClient(client)
		conn.Close()
		Manager.Broadcast(roomCode, Message{
			Type: "PLAYER_LEFT",
			Payload: map[string]string{
				"playerId": client.PlayerID,
			},
		})
	}()
	for {
		var message Message
		err := conn.ReadJSON(&message)
		if err != nil {
			break
		}
		log.Printf("Received message type: %s with payload: %+v", message.Type, message.Payload)

		switch message.Type {
		case "JOIN":
			if payload, ok := message.Payload.(map[string]interface{}); ok {
				if playerId, ok := payload["playerId"].(string); ok {
					client.PlayerID = playerId
					Manager.Broadcast(roomCode, Message{
						Type:    "PLAYER_JOINED",
						Payload: payload,
					})
				}
			}
		case "TRANSFER":
			if err := Manager.handleTransfer(client, message); err != nil {
				log.Printf("Transfer error: %v", err)
				conn.WriteJSON(Message{
					Type:    "ERROR",
					Payload: err.Error(),
				})
			} else {
				log.Printf("Transfer successful")
			}
		}
	}
}
