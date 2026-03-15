package middleware

import (
	"github.com/gin-gonic/gin"
)

// DevCORSMiddleware is a permissive CORS middleware for development
// Allows all origins with common headers and methods
func DevCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// For development, allow all origins
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")  // Wildcard for development
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "*")
			c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH, HEAD")
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}