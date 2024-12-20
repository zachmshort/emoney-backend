package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zachmshort/monopoly-backend/config"
	"github.com/zachmshort/monopoly-backend/routes"
)

func main() {

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://emoney.club",
			"https://www.emoney.club",
			"ws://localhost:3000",
			"wss://emoney.club",
			"wss://www.emoney.club",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"environment": map[string]bool{
				"DATABASE_URL":  os.Getenv("DATABASE_URL") != "",
				"DATABASE_NAME": os.Getenv("DATABASE_NAME") != "",
			},
		})
	})

	config.ConnectDB()
	routes.PropertyRoutes(r)
	routes.RoomRoutes(r)
	routes.PlayerRoutes(r)
	routes.TransferRoutes(r)
	routes.WebSocketRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
