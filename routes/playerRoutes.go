package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zachmshort/monopoly-backend/controllers"
)

func PlayerRoutes(r *gin.Engine) {
	player := r.Group("/player")
	{
		player.GET("/room/:roomCode", controllers.GetPlayersInRoom)
		player.GET("/:playerId/details", controllers.GetPlayerDetails)
	}
}
