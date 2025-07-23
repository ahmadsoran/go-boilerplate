// internal/middleware/logging.go
package middleware

import (
	"time"

	"your_project/internal/logger"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs details of each request
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process the request
		c.Next()

		// Calculate response time
		duration := time.Since(start)

		// Get the request ID from the context
		requestID := GetRequestID(c.Request.Context())

		// Log the request details with request ID
		logger.APILog.Infow("Request received",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", duration,
			"client_ip", c.ClientIP(),
		)
	}
}
