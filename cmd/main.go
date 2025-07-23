// cmd/main.go
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"your_project/configs"
	"your_project/internal/api"
	"your_project/internal/db"
	"your_project/internal/initializer"
	"your_project/internal/logger"
	"your_project/migrations"
)

func main() {
	logger.Init()

	// Load environment variables from .env file (for local development)
	// This will not overwrite existing environment variables
	if err := godotenv.Load(); err != nil {
		// Log a warning if the .env file is not found, but don't exit
		logger.Log.Warnw("Error loading .env file, using environment variables", "error", err)
	}

	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		logger.Log.Fatalw("Failed to load configuration", "error", err)
	}

	// Initialize database connection
	dbConn, err := db.Init(config.DatabaseURL)
	if err != nil {
		logger.Log.Fatalw("DB connection failed", "error", err)
	}

	// Run Go-based migrations if AUTO_MIGRATE is true
	if config.AutoMigrate {
		logger.Log.Infow("Running Go-based migrations...")
		if err := migrations.AutoMigrate(dbConn); err != nil {
			logger.Log.Fatalw("Go-based migration failed", "error", err)
		}
		logger.Log.Infow("Go-based migrations applied successfully.")
	}

	// Initialize repositories, services, and handlers using the initializer pattern
	repos := initializer.NewRepositoryContainer(dbConn)
	services := initializer.NewServiceContainer(repos, dbConn)
	handlers := initializer.NewHandlerContainer(services, dbConn)

	// Set up Gin router
	r := gin.Default()

	// Setup routes and apply middleware
	api.SetupRoutes(r, handlers)

	// Create the HTTP server
	serverAddr := fmt.Sprintf(":%s", config.Port)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	// Run the server in a goroutine
	go func() {
		logger.Log.Infof("Server started on port %s", config.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatalw("Server failed to start", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL (cannot be caught or ignored)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Infow("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatalw("Server forced to shutdown", "error", err)
	}

	logger.Log.Infow("Server exiting")
}

// Remove the runMigrations function that used golang-migrate/migrate
