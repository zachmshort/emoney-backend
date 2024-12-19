package routes

import (
	"os"

	"github.com/gin-gonic/gin"
)

func SystemRoutes(r *gin.Engine) {
	r.GET("/system/check", func(c *gin.Context) {
		envVars := []string{"DATABASE_URL", "DATABASE_NAME", "JWT_SECRET"}
		status := make(map[string]bool)

		for _, env := range envVars {
			status[env] = os.Getenv(env) != ""
		}

		c.JSON(200, gin.H{
			"status":       "ok",
			"environment":  status,
			"isProduction": os.Getenv("GIN_MODE") == "release",
		})
	})
}
