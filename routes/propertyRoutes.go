package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zachmshort/monopoly-backend/controllers"
)

func PropertyRoutes(r *gin.Engine) {
	property := r.Group("/property")
	{
		property.POST("/:id/player/:playerId", controllers.AddProperty)
		property.DELETE("/:id/player/:playerId", controllers.RemoveProperty)
		property.POST("/:id/player/:playerId/mortgage", controllers.MortgageProperty)
	}
}
