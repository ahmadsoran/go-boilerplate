package handlers

import (
	"net/http"

	"your_project/internal/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	*BaseHandler
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		BaseHandler: NewBaseHandler(),
		db:          db,
	}
}

func (h *HealthHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.CheckHealth)
}

func (h *HealthHandler) CheckHealth(c *gin.Context) {
	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		h.ErrorHandler.HandleError(c, pkg.NewDatabaseConnectionError("postgresql", err, "Failed to get database instance"))
		return
	}
	if err := sqlDB.Ping(); err != nil {
		h.ErrorHandler.HandleError(c, pkg.NewDatabaseConnectionError("postgresql", err, "Database ping failed"))
		return
	}

	// If all checks pass, return OK
	c.JSON(http.StatusOK, gin.H{"status": "UP", "database": "UP"})
}
