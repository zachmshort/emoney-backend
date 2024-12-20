package websocket

import (
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
	if roomCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room code required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &Client{
		Conn: conn,
		Room: roomCode,
	}

	Manager.AddClient(client)
	defer func() {
		Manager.RemoveClient(client)
		conn.Close()
		Manager.Broadcast(roomCode, Message{
			Type: "PLAYER_LEFT",
			Payload: map[string]string{
				"deviceId": client.DeviceID,
			},
		})
	}()

	for {
		var message Message
		err := conn.ReadJSON(&message)
		if err != nil {
			break
		}

		switch message.Type {
		case "JOIN":
			if payload, ok := message.Payload.(map[string]interface{}); ok {
				if deviceId, ok := payload["deviceId"].(string); ok {
					client.DeviceID = deviceId
					Manager.Broadcast(roomCode, Message{
						Type:    "PLAYER_JOINED",
						Payload: payload,
					})
				}
			}
		case "TRANSFER":
			if err := Manager.handleTransfer(client, message); err != nil {
				conn.WriteJSON(Message{
					Type:    "ERROR",
					Payload: err.Error(),
				})
			}
		}
	}
}
