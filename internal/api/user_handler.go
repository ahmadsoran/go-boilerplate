// internal/api/user_handler.go
package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"your_project/internal/logger"
	"your_project/internal/model"
	"your_project/internal/pkg"
	"your_project/internal/service"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.POST("", h.CreateUser)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteItem)
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the user struct
	if err := validate.Struct(user); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors.Error()})
		return
	}

	// Pass the request context to the service layer
	if err := h.svc.CreateUser(c.Request.Context(), &user); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.handleError(c, pkg.NewInvalidInputError("invalid user ID"))
		return
	}

	// Pass the request context to the service layer
	user, err := h.svc.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		h.handleError(c, err)
		return	
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.handleError(c, pkg.NewInvalidInputError("invalid user ID"))
		return
	}

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the user struct
	if err := validate.Struct(user); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors.Error()})
		return
	}

	user.ID = uint(id)

	// Pass the request context to the service layer
	if err := h.svc.UpdateUser(c.Request.Context(), &user); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.handleError(c, pkg.NewInvalidInputError("invalid user ID"))
		return
	}

	// Pass the request context to the service layer
	if err := h.svc.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// handleError maps custom errors to HTTP status codes and returns a JSON response
func (h *UserHandler) handleError(c *gin.Context, err error) {
	var notFoundErr *pkg.NotFoundError
	var invalidInputErr *pkg.InvalidInputError
	var internalServerErr *pkg.InternalServerError

	switch {
	case errors.As(err, &notFoundErr):
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
	case errors.As(err, &invalidInputErr):
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidInputErr.Error()})
	case errors.As(err, &internalServerErr):
		// Log the original error for internal server errors
		logger.Log.Errorw("Internal Server Error", "error", internalServerErr.Err, "message", internalServerErr.Message)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	default:
		// Log unexpected errors
		logger.Log.Errorw("Unexpected Error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
}
