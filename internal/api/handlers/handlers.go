package handlers

import (
	"your_project/internal/logger"
	"your_project/internal/pkg"
)

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	ErrorHandler *pkg.HTTPErrorHandler
}

// NewBaseHandler creates a new base handler with error handling
func NewBaseHandler() *BaseHandler {
	return &BaseHandler{
		ErrorHandler: pkg.NewHTTPErrorHandler(logger.APILog),
	}
}
