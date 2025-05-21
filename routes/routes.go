package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zachmshort/emoney-backend/controllers"
	"github.com/zachmshort/emoney-backend/websocket"
)

func Routes(r *gin.Engine) {
	r.GET("/ws/room/:code", websocket.HandleWebSocket)

	room := r.Group("/room")
	{
		room.POST("", controllers.CreateRoom)
		room.POST("/join", controllers.JoinRoom)

		room.GET("/:code/players", controllers.GetPlayersInRoom)

		room.GET("/:code/properties", controllers.GetAvailableProperties)
	}

	player := r.Group("/player")
	{
		player.GET("/:playerId/details", controllers.GetPlayerDetails)

		playerProperty := player.Group("/:playerId/property/:propertyId")
		{
			playerProperty.POST("", controllers.AddProperty)
			playerProperty.DELETE("", controllers.RemoveProperty)
			playerProperty.POST("/mortgage", controllers.MortgageProperty)
		}
	}
}
