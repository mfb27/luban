package middleware

import (
	"github.com/gin-gonic/gin"
)

// SimpleCORSMiddleware is a very simple CORS middleware for development
func SimpleCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add CORS headers for all requests
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		// c.Header("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
