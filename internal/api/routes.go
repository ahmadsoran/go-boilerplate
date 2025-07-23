package api

import (
	"your_project/configs"
	"your_project/internal/initializer"
	"your_project/internal/middleware"
	"your_project/internal/pkg"

	"github.com/gin-gonic/gin"
)

// SetupRoutes registers all API routes and applies middleware
func SetupRoutes(r *gin.Engine, handlers *initializer.HandlerContainer, config configs.Config) {
	// Initialize JWT manager for middleware
	jwtManager := pkg.NewJWTManager(config.JWTSecret, config.JWTExpiryHours)

	// Apply global middleware
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.LoggingMiddleware())

	// Register health check route (no middleware needed for health check)
	handlers.Health.RegisterRoutes(r)

	// Group routes by functionality or version
	apiRoutes := r.Group("/api")
	{
		// Authentication routes (public) - handled within the user handler
		authRoutes := apiRoutes.Group("/auth")
		{
			authRoutes.POST("/signup", handlers.User.SignUp)
			authRoutes.POST("/login", handlers.User.Login)
			authRoutes.POST("/refresh", handlers.User.RefreshToken)
		}

		// Protected user routes
		protectedUsers := apiRoutes.Group("/users")
		protectedUsers.Use(middleware.AuthMiddleware(jwtManager))
		{
			protectedUsers.GET("/:id", handlers.User.GetUser)
			protectedUsers.PUT("/:id", handlers.User.UpdateUser)
			protectedUsers.DELETE("/:id", handlers.User.DeleteItem)
		}

		// Add other module routes here
	}
}
