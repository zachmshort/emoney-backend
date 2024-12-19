package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zachmshort/monopoly-backend/controllers"
)

func RoomRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/create", controllers.CreateRoom)
		auth.POST("/join", controllers.JoinRoom)
	}

}
