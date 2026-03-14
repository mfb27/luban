package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestIDHeader is the header key for request ID
const RequestIDHeader = "X-Request-ID"

// WithRequestID creates a middleware that generates a unique request ID for each request
// and adds it to the context for logging
func WithRequestID(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID is already provided in the header
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			// Generate a new UUID if not provided
			requestID = uuid.New().String()
		}

		// Set request ID in the context
		c.Set("request_id", requestID)

		// Add request ID to response headers
		c.Header(RequestIDHeader, requestID)

		// Add request ID to logger for this request
		log = log.With(zap.String("request_id", requestID))
		c.Set("logger", log)

		// Continue with the next handler
		c.Next()
	}
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return ""
}

// GetLoggerWithRequestID retrieves the logger with request ID from the context
func GetLoggerWithRequestID(c *gin.Context) *zap.Logger {
	if logger, exists := c.Get("logger"); exists {
		return logger.(*zap.Logger)
	}
	return zap.L()
}

// GetRequestIDField returns a zap.Field with the request ID
func GetRequestIDField(c *gin.Context) zap.Field {
	return zap.String("request_id", GetRequestID(c))
}

// GetErrorField returns a zap.Field with the error
func GetErrorField(err error) zap.Field {
	return zap.Error(err)
}

// GetTraceField returns a zap.Field for tracing
func GetTraceField(traceID string) zap.Field {
	return zap.String("trace_id", traceID)
}
