package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zachmshort/emoney-backend/controllers"
	"github.com/zachmshort/emoney-backend/websocket"
)

func Routes(r *gin.Engine) {
	apiVersion := r.Group("/v1")
	apiVersion.GET("/ws/room/:code", websocket.HandleWebSocket)

	rooms := apiVersion.Group("/rooms")
	{
		rooms.POST("", controllers.CreateRoom)

		room := rooms.Group("/:code")
		room.GET("/players", controllers.GetPlayersInRoom)
		room.GET("/properties", controllers.GetAvailableProperties)
		room.GET("/exists", controllers.CheckIfRoomCodeExists)

		players := room.Group("/players")
		{
			players.POST("", controllers.JoinRoom)
			player := players.Group("/:playerId")

			{
				player.GET("", controllers.GetPlayerDetails)

				properties := player.Group("/properties")
				{
					property := properties.Group("/:propertyId")

					property.POST("/mortgage", controllers.MortgageProperty)
					property.POST("", controllers.AddProperty)
					property.DELETE("", controllers.RemoveProperty)
				}
			}
		}
	}
}
