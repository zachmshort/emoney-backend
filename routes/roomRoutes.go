package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zachmshort/monopoly-backend/controllers"
)

func RoomRoutes(r *gin.Engine) {
	r.POST("/room", controllers.CreateRoom)
	r.POST("/room/join", controllers.JoinRoom)
	r.GET("/ws/:roomCode", controllers.HandleWebSocket)
}
