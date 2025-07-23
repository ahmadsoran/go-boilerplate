// internal/api/routes.go
package api

import (
	"github.com/gin-gonic/gin"
	"your_project/internal/initializer"
	"your_project/internal/middleware"
)

// SetupRoutes registers all API routes and applies middleware
func SetupRoutes(r *gin.Engine, handlers *initializer.HandlerContainer) {
	// Apply global middleware
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.LoggingMiddleware())

	// Register health check route (no middleware needed for health check)
	handlers.Health.RegisterRoutes(r)

	// Group routes by functionality or version
	apiRoutes := r.Group("/api")
	{
		// User routes
		handlers.User.RegisterRoutes(apiRoutes)

		// Add other module routes here (e.g., handlers.Product.RegisterRoutes(apiRoutes))
	}
}
