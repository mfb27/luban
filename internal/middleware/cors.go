package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	// Check if we're in development mode
	isDevelopment := os.Getenv("GIN_MODE") != "release"

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin != "" {
			// In development, allow all local origins
			if isDevelopment && (strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1")) {
				// Allow all local development origins
				c.Header("Access-Control-Allow-Origin", origin)
			} else {
				// In production, check against allowed origins
				allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
				if len(allowedOrigins) == 1 && allowedOrigins[0] == "" {
					// Default: allow same origin
					c.Header("Access-Control-Allow-Origin", origin)
				} else {
					for _, allowedOrigin := range allowedOrigins {
						if strings.TrimSpace(allowedOrigin) == origin {
							c.Header("Access-Control-Allow-Origin", origin)
							break
						}
					}
				}
			}

			// Common CORS headers
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Api-Key")
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