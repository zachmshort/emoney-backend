package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zachmshort/monopoly-backend/controllers"
)

func TransferRoutes(r *gin.Engine) {
	transaction := r.Group("/transactions")
	{
		transaction.POST("/bank/:playerId", controllers.BankTransfer)
		transaction.POST("/request", controllers.RequestTransfer)
		transaction.POST("/transfer", controllers.Transfer)
	}
}
