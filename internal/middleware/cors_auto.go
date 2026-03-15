package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

// AutoCORSMiddleware automatically chooses the appropriate CORS middleware based on environment
func AutoCORSMiddleware() gin.HandlerFunc {
	// Check if we're in development mode
	// GIN_MODE=debug or not set
	ginMode := os.Getenv("GIN_MODE")
	isDevelopment := ginMode == "" || ginMode == "debug"

	// Also check if we're explicitly in dev mode
	isDev := os.Getenv("APP_ENV") == "development" ||
		os.Getenv("GO_ENV") == "development" ||
		os.Getenv("NODE_ENV") == "development"

	if isDevelopment || isDev {
		// Use development CORS - more permissive
		return DevCORSMiddleware()
	}

	// Use production CORS - more restrictive
	return CORSMiddleware()
}