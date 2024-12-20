package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zachmshort/monopoly-backend/websocket"
)

func WebSocketRoutes(r *gin.Engine) {
	r.GET("/ws/room/:code", websocket.HandleWebSocket)
}
