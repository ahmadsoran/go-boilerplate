// internal/middleware/requestid.go
package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const (
	requestIDKey contextKey = "requestID"
)

// RequestIDMiddleware adds a unique request ID to the context and response headers
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()

		// Add the request ID to the context
		ctx := context.WithValue(c.Request.Context(), requestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		// Add the request ID to the response headers
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}

// GetRequestID returns the request ID from the context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}
	return ""
}
