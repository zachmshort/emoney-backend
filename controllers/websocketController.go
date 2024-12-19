// controllers/wsController.go
package controllers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type RoomManager struct {
	rooms map[string]map[*websocket.Conn]bool
	mu    sync.RWMutex
}

var manager = RoomManager{
	rooms: make(map[string]map[*websocket.Conn]bool),
}

func HandleWebSocket(c *gin.Context) {
	roomCode := c.Param("roomCode")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	manager.mu.Lock()
	if _, exists := manager.rooms[roomCode]; !exists {
		manager.rooms[roomCode] = make(map[*websocket.Conn]bool)
	}
	manager.rooms[roomCode][conn] = true
	manager.mu.Unlock()

	defer func() {
		manager.mu.Lock()
		delete(manager.rooms[roomCode], conn)
		if len(manager.rooms[roomCode]) == 0 {
			delete(manager.rooms, roomCode)
		}
		manager.mu.Unlock()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		HandleMessage(roomCode, message)
	}
}

func HandleMessage(roomCode string, message []byte) {
}

func BroadcastToRoom(roomCode string, message interface{}) {
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	if connections, exists := manager.rooms[roomCode]; exists {
		for conn := range connections {
			err := conn.WriteJSON(message)
			if err != nil {
				conn.Close()
				delete(connections, conn)
			}
		}
	}
}
