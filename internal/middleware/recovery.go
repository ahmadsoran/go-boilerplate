// internal/middleware/recovery.go
package middleware

import (
	"net/http"
	"runtime/debug"

	"your_project/internal/logger"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware recovers from panics and returns an Internal Server Error
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace and request ID
				requestID := GetRequestID(c.Request.Context())
				logger.APILog.Errorw("Panic recovered",
					"request_id", requestID,
					"panic", err,
					"stack", string(debug.Stack()),
				)

				// Abort the request and return an Internal Server Error
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":      "Internal Server Error",
					"request_id": requestID, // Include request ID in the error response
				})
			}
		}()

		c.Next()
	}
}
