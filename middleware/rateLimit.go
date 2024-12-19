package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimitMiddleware(rps float64) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(rps), int(rps))
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(429, gin.H{
				"error": "too many requests",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
