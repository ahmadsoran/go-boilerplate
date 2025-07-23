package handlers

import (
	"net/http"
	"strconv"
	"time"

	"your_project/configs"
	"your_project/internal/model"
	"your_project/internal/pkg"
	"your_project/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type UserHandler struct {
	*BaseHandler
	svc        service.UserService
	jwtManager *pkg.JWTManager
}

func NewUserHandler(svc service.UserService, config configs.Config) *UserHandler {
	jwtManager := pkg.NewJWTManager(config.JWTSecret, config.JWTExpiryHours)
	return &UserHandler{
		BaseHandler: NewBaseHandler(),
		svc:         svc,
		jwtManager:  jwtManager,
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
		h.ErrorHandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.ErrorHandler.HandleError(c, pkg.NewInvalidInputError("invalid user ID"))
		return
	}

	// Pass the request context to the service layer
	user, err := h.svc.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		h.ErrorHandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.ErrorHandler.HandleError(c, pkg.NewInvalidInputError("invalid user ID"))
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
		h.ErrorHandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.ErrorHandler.HandleError(c, pkg.NewInvalidInputError("invalid user ID"))
		return
	}

	// Pass the request context to the service layer
	if err := h.svc.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		h.ErrorHandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// SignUp User
func (h *UserHandler) SignUp(c *gin.Context) {
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
	if err := h.svc.RegisterUser(c.Request.Context(), &user); err != nil {
		h.ErrorHandler.HandleError(c, err)
		return
	}
	// Optionally, you can generate a token here or just return success
	// For simplicity, we will just return a success message
	// You can also generate a token pair here if needed
	tokenPair, err := h.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		h.ErrorHandler.HandleError(c, pkg.NewInternalServerError(err, "Failed to generate tokens"))
		return
	}
	// Store refresh token in database
	refreshExpiry := time.Now().Add(time.Duration(h.jwtManager.RefreshExpiryHours()) * time.Hour)
	err = h.svc.UpdateRefreshToken(c.Request.Context(), user.ID, tokenPair.RefreshToken, refreshExpiry)
	if err != nil {
		h.ErrorHandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
func (h *UserHandler) Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.LoginUser(c.Request.Context(), loginData.Email, loginData.Password)
	if err != nil {
		h.ErrorHandler.HandleError(c, err)
		return
	}

	// Generate token pair
	tokenPair, err := h.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		h.ErrorHandler.HandleError(c, pkg.NewInternalServerError(err, "Failed to generate tokens"))
		return
	}

	// Store refresh token in database
	refreshExpiry := time.Now().Add(time.Duration(h.jwtManager.RefreshExpiryHours()) * time.Hour)
	err = h.svc.UpdateRefreshToken(c.Request.Context(), user.ID, tokenPair.RefreshToken, refreshExpiry)
	if err != nil {
		h.ErrorHandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"user":          user,
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
	})
}

// RefreshToken handles refresh token requests
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var refreshData struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&refreshData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the refresh token
	claims, err := h.jwtManager.ValidateToken(refreshData.RefreshToken)
	if err != nil {
		h.ErrorHandler.HandleError(c, pkg.NewUnauthorizedError("Invalid refresh token"))
		return
	}

	// Check if it's actually a refresh token
	if claims.TokenType != "refresh" {
		h.ErrorHandler.HandleError(c, pkg.NewUnauthorizedError("Invalid token type - refresh token required"))
		return
	}

	// Get user by refresh token to ensure it exists in database
	user, err := h.svc.RefreshToken(c.Request.Context(), refreshData.RefreshToken)
	if err != nil {
		h.ErrorHandler.HandleError(c, err)
		return
	}

	// Generate new token pair
	tokenPair, err := h.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		h.ErrorHandler.HandleError(c, pkg.NewInternalServerError(err, "Failed to generate new tokens"))
		return
	}

	// Update refresh token in database
	refreshExpiry := time.Now().Add(time.Duration(h.jwtManager.RefreshExpiryHours()) * time.Hour)
	err = h.svc.UpdateRefreshToken(c.Request.Context(), user.ID, tokenPair.RefreshToken, refreshExpiry)
	if err != nil {
		h.ErrorHandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Token refreshed successfully",
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
	})
}
