package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// responseBodyWriter wraps gin.ResponseWriter to capture response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// RequestLogger logs request and response details
func RequestLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for static assets and health check
		if c.Request.URL.Path == "/api/health" ||
			strings.HasPrefix(c.Request.URL.Path, "/assets") ||
			strings.HasPrefix(c.Request.URL.Path, "/favicon.ico") {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Get logger with request ID
		logger := GetLoggerWithRequestID(c)

		// Capture request body
		var requestBody []byte
		if c.Request.Body != nil && c.Request.Method != "GET" {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Wrap response writer to capture response body
		writer := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = writer

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Prepare log fields
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		// Add query parameters
		if query != "" {
			fields = append(fields, zap.String("query", query))
		}

		// Add request body (truncate if too large)
		if len(requestBody) > 0 {
			bodyStr := string(requestBody)
			if len(bodyStr) > 1000 {
				bodyStr = bodyStr[:1000] + "... (truncated)"
			}
			fields = append(fields, zap.String("request_body", bodyStr))
		}

		// Add response body (truncate if too large)
		responseBody := writer.body.Bytes()
		if len(responseBody) > 0 {
			bodyStr := string(responseBody)
			if len(bodyStr) > 1000 {
				bodyStr = bodyStr[:1000] + "... (truncated)"
			}
			fields = append(fields, zap.String("response_body", bodyStr))
		}

		// Log based on status code
		if c.Writer.Status() >= 500 {
			logger.Error("HTTP Request", fields...)
		} else if c.Writer.Status() >= 400 {
			logger.Warn("HTTP Request", fields...)
		} else {
			logger.Info("HTTP Request", fields...)
		}
	}
}
