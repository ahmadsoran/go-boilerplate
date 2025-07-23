// internal/pkg/errors.go
package pkg

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Custom error types

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewNotFoundError(format string, a ...interface{}) error {
	return &NotFoundError{Message: fmt.Sprintf(format, a...)}
}

type InvalidInputError struct {
	Message string
}

func (e *InvalidInputError) Error() string {
	return e.Message
}

func NewInvalidInputError(format string, a ...interface{}) error {
	return &InvalidInputError{Message: fmt.Sprintf(format, a...)}
}

type InternalServerError struct {
	Message string
	Err     error // Original error
}

func (e *InternalServerError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewInternalServerError(err error, format string, a ...interface{}) error {
	return &InternalServerError{Message: fmt.Sprintf(format, a...), Err: err}
}

// DuplicateError represents a duplicate entry error (e.g., unique constraint violation)
type DuplicateError struct {
	Message string
	Field   string // The field that caused the duplicate error
}

func (e *DuplicateError) Error() string {
	return e.Message
}

func NewDuplicateError(field, format string, a ...interface{}) error {
	return &DuplicateError{
		Message: fmt.Sprintf(format, a...),
		Field:   field,
	}
}

// ValidationError represents validation errors with field-specific details
type ValidationError struct {
	Message string
	Field   string
	Value   interface{}
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(field string, value interface{}, format string, a ...interface{}) error {
	return &ValidationError{
		Message: fmt.Sprintf(format, a...),
		Field:   field,
		Value:   value,
	}
}

// UnauthorizedError represents authentication failures
type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

func NewUnauthorizedError(format string, a ...interface{}) error {
	return &UnauthorizedError{Message: fmt.Sprintf(format, a...)}
}

// ForbiddenError represents authorization failures (user authenticated but lacks permission)
type ForbiddenError struct {
	Message  string
	Resource string // The resource they tried to access
	Action   string // The action they tried to perform
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

func NewForbiddenError(resource, action, format string, a ...interface{}) error {
	return &ForbiddenError{
		Message:  fmt.Sprintf(format, a...),
		Resource: resource,
		Action:   action,
	}
}

// ConflictError represents business logic conflicts
type ConflictError struct {
	Message string
	Details string
}

func (e *ConflictError) Error() string {
	return e.Message
}

func NewConflictError(details, format string, a ...interface{}) error {
	return &ConflictError{
		Message: fmt.Sprintf(format, a...),
		Details: details,
	}
}

// RateLimitError represents rate limiting errors
type RateLimitError struct {
	Message   string
	RetryTime int // Seconds until retry is allowed
}

func (e *RateLimitError) Error() string {
	return e.Message
}

func NewRateLimitError(retryTime int, format string, a ...interface{}) error {
	return &RateLimitError{
		Message:   fmt.Sprintf(format, a...),
		RetryTime: retryTime,
	}
}

// ServiceUnavailableError represents external service failures
type ServiceUnavailableError struct {
	Message     string
	ServiceName string
	Err         error
}

func (e *ServiceUnavailableError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewServiceUnavailableError(serviceName string, err error, format string, a ...interface{}) error {
	return &ServiceUnavailableError{
		Message:     fmt.Sprintf(format, a...),
		ServiceName: serviceName,
		Err:         err,
	}
}

// TimeoutError represents request timeout errors
type TimeoutError struct {
	Message     string
	TimeoutSecs int
	Operation   string
}

func (e *TimeoutError) Error() string {
	return e.Message
}

func NewTimeoutError(operation string, timeoutSecs int, format string, a ...interface{}) error {
	return &TimeoutError{
		Message:     fmt.Sprintf(format, a...),
		TimeoutSecs: timeoutSecs,
		Operation:   operation,
	}
}

// PayloadTooLargeError represents request payload size errors
type PayloadTooLargeError struct {
	Message    string
	MaxSize    int64
	ActualSize int64
}

func (e *PayloadTooLargeError) Error() string {
	return e.Message
}

func NewPayloadTooLargeError(maxSize, actualSize int64, format string, a ...interface{}) error {
	return &PayloadTooLargeError{
		Message:    fmt.Sprintf(format, a...),
		MaxSize:    maxSize,
		ActualSize: actualSize,
	}
}

// UnsupportedMediaTypeError represents unsupported content type errors
type UnsupportedMediaTypeError struct {
	Message        string
	ReceivedType   string
	SupportedTypes []string
}

func (e *UnsupportedMediaTypeError) Error() string {
	return e.Message
}

func NewUnsupportedMediaTypeError(receivedType string, supportedTypes []string, format string, a ...interface{}) error {
	return &UnsupportedMediaTypeError{
		Message:        fmt.Sprintf(format, a...),
		ReceivedType:   receivedType,
		SupportedTypes: supportedTypes,
	}
}

// BadGatewayError represents upstream service errors
type BadGatewayError struct {
	Message         string
	UpstreamService string
	Err             error
}

func (e *BadGatewayError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewBadGatewayError(upstreamService string, err error, format string, a ...interface{}) error {
	return &BadGatewayError{
		Message:         fmt.Sprintf(format, a...),
		UpstreamService: upstreamService,
		Err:             err,
	}
}

// TooManyRequestsError represents rate limiting by external services
type TooManyRequestsError struct {
	Message    string
	Service    string
	RetryAfter int // Seconds
}

func (e *TooManyRequestsError) Error() string {
	return e.Message
}

func NewTooManyRequestsError(service string, retryAfter int, format string, a ...interface{}) error {
	return &TooManyRequestsError{
		Message:    fmt.Sprintf(format, a...),
		Service:    service,
		RetryAfter: retryAfter,
	}
}

// DatabaseConnectionError represents database connectivity issues
type DatabaseConnectionError struct {
	Message  string
	Database string
	Err      error
}

func (e *DatabaseConnectionError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewDatabaseConnectionError(database string, err error, format string, a ...interface{}) error {
	return &DatabaseConnectionError{
		Message:  fmt.Sprintf(format, a...),
		Database: database,
		Err:      err,
	}
}

// MigrationError represents database migration failures
type MigrationError struct {
	Message       string
	MigrationName string
	Err           error
}

func (e *MigrationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewMigrationError(migrationName string, err error, format string, a ...interface{}) error {
	return &MigrationError{
		Message:       fmt.Sprintf(format, a...),
		MigrationName: migrationName,
		Err:           err,
	}
}

// ConfigurationError represents configuration-related errors
type ConfigurationError struct {
	Message      string
	ConfigKey    string
	ExpectedType string
}

func (e *ConfigurationError) Error() string {
	return e.Message
}

func NewConfigurationError(configKey, expectedType, format string, a ...interface{}) error {
	return &ConfigurationError{
		Message:      fmt.Sprintf(format, a...),
		ConfigKey:    configKey,
		ExpectedType: expectedType,
	}
}

// FileNotFoundError represents file system errors
type FileNotFoundError struct {
	Message  string
	FilePath string
}

func (e *FileNotFoundError) Error() string {
	return e.Message
}

func NewFileNotFoundError(filePath, format string, a ...interface{}) error {
	return &FileNotFoundError{
		Message:  fmt.Sprintf(format, a...),
		FilePath: filePath,
	}
}

// PermissionDeniedError represents file system permission errors
type PermissionDeniedError struct {
	Message   string
	Resource  string
	Operation string
}

func (e *PermissionDeniedError) Error() string {
	return e.Message
}

func NewPermissionDeniedError(resource, operation, format string, a ...interface{}) error {
	return &PermissionDeniedError{
		Message:   fmt.Sprintf(format, a...),
		Resource:  resource,
		Operation: operation,
	}
}

// NetworkError represents network connectivity issues
type NetworkError struct {
	Message string
	Host    string
	Port    int
	Err     error
}

func (e *NetworkError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewNetworkError(host string, port int, err error, format string, a ...interface{}) error {
	return &NetworkError{
		Message: fmt.Sprintf(format, a...),
		Host:    host,
		Port:    port,
		Err:     err,
	}
}

// CacheError represents caching system errors
type CacheError struct {
	Message   string
	CacheType string // Redis, Memcached, etc.
	Key       string
	Err       error
}

func (e *CacheError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewCacheError(cacheType, key string, err error, format string, a ...interface{}) error {
	return &CacheError{
		Message:   fmt.Sprintf(format, a...),
		CacheType: cacheType,
		Key:       key,
		Err:       err,
	}
}

// QueueError represents message queue errors
type QueueError struct {
	Message   string
	QueueName string
	Operation string // publish, consume, acknowledge, etc.
	Err       error
}

func (e *QueueError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewQueueError(queueName, operation string, err error, format string, a ...interface{}) error {
	return &QueueError{
		Message:   fmt.Sprintf(format, a...),
		QueueName: queueName,
		Operation: operation,
		Err:       err,
	}
}

// ExternalAPIError represents third-party API errors
type ExternalAPIError struct {
	Message    string
	APIName    string
	StatusCode int
	Endpoint   string
	Err        error
}

func (e *ExternalAPIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewExternalAPIError(apiName, endpoint string, statusCode int, err error, format string, a ...interface{}) error {
	return &ExternalAPIError{
		Message:    fmt.Sprintf(format, a...),
		APIName:    apiName,
		StatusCode: statusCode,
		Endpoint:   endpoint,
		Err:        err,
	}
}

// HTTPErrorHandler provides centralized error handling for HTTP responses
type HTTPErrorHandler struct {
	logger Logger // Interface for logging
}

// Logger interface to avoid tight coupling with specific logging library
type Logger interface {
	Errorw(msg string, keysAndValues ...interface{})
}

// NewHTTPErrorHandler creates a new HTTP error handler
func NewHTTPErrorHandler(logger Logger) *HTTPErrorHandler {
	return &HTTPErrorHandler{logger: logger}
}

// HandleError maps custom errors to HTTP status codes and returns a JSON response
func (h *HTTPErrorHandler) HandleError(c *gin.Context, err error) {
	var notFoundErr *NotFoundError
	var invalidInputErr *InvalidInputError
	var internalServerErr *InternalServerError
	var duplicateErr *DuplicateError
	var validationErr *ValidationError
	var unauthorizedErr *UnauthorizedError
	var forbiddenErr *ForbiddenError
	var conflictErr *ConflictError
	var rateLimitErr *RateLimitError
	var serviceUnavailableErr *ServiceUnavailableError
	var timeoutErr *TimeoutError
	var payloadTooLargeErr *PayloadTooLargeError
	var unsupportedMediaTypeErr *UnsupportedMediaTypeError
	var badGatewayErr *BadGatewayError
	var tooManyRequestsErr *TooManyRequestsError
	var databaseConnectionErr *DatabaseConnectionError
	var migrationErr *MigrationError
	var configurationErr *ConfigurationError
	var fileNotFoundErr *FileNotFoundError
	var permissionDeniedErr *PermissionDeniedError
	var networkErr *NetworkError
	var cacheErr *CacheError
	var queueErr *QueueError
	var externalAPIErr *ExternalAPIError

	switch {
	case errors.As(err, &notFoundErr):
		c.JSON(http.StatusNotFound, gin.H{
			"error": notFoundErr.Error(),
			"type":  "not_found",
		})

	case errors.As(err, &invalidInputErr):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": invalidInputErr.Error(),
			"type":  "invalid_input",
		})

	case errors.As(err, &duplicateErr):
		c.JSON(http.StatusConflict, gin.H{
			"error": duplicateErr.Error(),
			"type":  "duplicate",
			"field": duplicateErr.Field,
		})

	case errors.As(err, &validationErr):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validationErr.Error(),
			"type":  "validation",
			"field": validationErr.Field,
			"value": validationErr.Value,
		})

	case errors.As(err, &unauthorizedErr):
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": unauthorizedErr.Error(),
			"type":  "unauthorized",
		})

	case errors.As(err, &forbiddenErr):
		c.JSON(http.StatusForbidden, gin.H{
			"error":    forbiddenErr.Error(),
			"type":     "forbidden",
			"resource": forbiddenErr.Resource,
			"action":   forbiddenErr.Action,
		})

	case errors.As(err, &conflictErr):
		c.JSON(http.StatusConflict, gin.H{
			"error":   conflictErr.Error(),
			"type":    "conflict",
			"details": conflictErr.Details,
		})

	case errors.As(err, &rateLimitErr):
		c.Header("Retry-After", fmt.Sprintf("%d", rateLimitErr.RetryTime))
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":       rateLimitErr.Error(),
			"type":        "rate_limit",
			"retry_after": rateLimitErr.RetryTime,
		})

	case errors.As(err, &serviceUnavailableErr):
		h.logger.Errorw("Service Unavailable", "service", serviceUnavailableErr.ServiceName, "error", serviceUnavailableErr.Err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Service temporarily unavailable",
			"type":    "service_unavailable",
			"service": serviceUnavailableErr.ServiceName,
		})

	case errors.As(err, &timeoutErr):
		h.logger.Errorw("Request Timeout", "operation", timeoutErr.Operation, "timeout", timeoutErr.TimeoutSecs)
		c.JSON(http.StatusRequestTimeout, gin.H{
			"error":     timeoutErr.Error(),
			"type":      "timeout",
			"operation": timeoutErr.Operation,
			"timeout":   timeoutErr.TimeoutSecs,
		})

	case errors.As(err, &payloadTooLargeErr):
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error":       payloadTooLargeErr.Error(),
			"type":        "payload_too_large",
			"max_size":    payloadTooLargeErr.MaxSize,
			"actual_size": payloadTooLargeErr.ActualSize,
		})

	case errors.As(err, &unsupportedMediaTypeErr):
		c.JSON(http.StatusUnsupportedMediaType, gin.H{
			"error":           unsupportedMediaTypeErr.Error(),
			"type":            "unsupported_media_type",
			"received_type":   unsupportedMediaTypeErr.ReceivedType,
			"supported_types": unsupportedMediaTypeErr.SupportedTypes,
		})

	case errors.As(err, &badGatewayErr):
		h.logger.Errorw("Bad Gateway", "upstream", badGatewayErr.UpstreamService, "error", badGatewayErr.Err)
		c.JSON(http.StatusBadGateway, gin.H{
			"error":    "Upstream service error",
			"type":     "bad_gateway",
			"upstream": badGatewayErr.UpstreamService,
		})

	case errors.As(err, &tooManyRequestsErr):
		c.Header("Retry-After", fmt.Sprintf("%d", tooManyRequestsErr.RetryAfter))
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":       tooManyRequestsErr.Error(),
			"type":        "too_many_requests",
			"service":     tooManyRequestsErr.Service,
			"retry_after": tooManyRequestsErr.RetryAfter,
		})

	case errors.As(err, &databaseConnectionErr):
		h.logger.Errorw("Database Connection Error", "database", databaseConnectionErr.Database, "error", databaseConnectionErr.Err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":    "Database connection failed",
			"type":     "database_connection",
			"database": databaseConnectionErr.Database,
		})

	case errors.As(err, &migrationErr):
		h.logger.Errorw("Migration Error", "migration", migrationErr.MigrationName, "error", migrationErr.Err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Database migration failed",
			"type":      "migration",
			"migration": migrationErr.MigrationName,
		})

	case errors.As(err, &configurationErr):
		h.logger.Errorw("Configuration Error", "key", configurationErr.ConfigKey, "expected_type", configurationErr.ExpectedType)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         "Configuration error",
			"type":          "configuration",
			"config_key":    configurationErr.ConfigKey,
			"expected_type": configurationErr.ExpectedType,
		})

	case errors.As(err, &fileNotFoundErr):
		c.JSON(http.StatusNotFound, gin.H{
			"error":     fileNotFoundErr.Error(),
			"type":      "file_not_found",
			"file_path": fileNotFoundErr.FilePath,
		})

	case errors.As(err, &permissionDeniedErr):
		c.JSON(http.StatusForbidden, gin.H{
			"error":     permissionDeniedErr.Error(),
			"type":      "permission_denied",
			"resource":  permissionDeniedErr.Resource,
			"operation": permissionDeniedErr.Operation,
		})

	case errors.As(err, &networkErr):
		h.logger.Errorw("Network Error", "host", networkErr.Host, "port", networkErr.Port, "error", networkErr.Err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Network connectivity issue",
			"type":  "network",
			"host":  networkErr.Host,
			"port":  networkErr.Port,
		})

	case errors.As(err, &cacheErr):
		h.logger.Errorw("Cache Error", "cache_type", cacheErr.CacheType, "key", cacheErr.Key, "error", cacheErr.Err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":      "Cache service error",
			"type":       "cache",
			"cache_type": cacheErr.CacheType,
			"key":        cacheErr.Key,
		})

	case errors.As(err, &queueErr):
		h.logger.Errorw("Queue Error", "queue", queueErr.QueueName, "operation", queueErr.Operation, "error", queueErr.Err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":     "Message queue error",
			"type":      "queue",
			"queue":     queueErr.QueueName,
			"operation": queueErr.Operation,
		})

	case errors.As(err, &externalAPIErr):
		h.logger.Errorw("External API Error", "api", externalAPIErr.APIName, "endpoint", externalAPIErr.Endpoint, "status", externalAPIErr.StatusCode, "error", externalAPIErr.Err)
		c.JSON(http.StatusBadGateway, gin.H{
			"error":       "External API error",
			"type":        "external_api",
			"api":         externalAPIErr.APIName,
			"endpoint":    externalAPIErr.Endpoint,
			"status_code": externalAPIErr.StatusCode,
		})

	case errors.As(err, &internalServerErr):
		// Log the original error for internal server errors
		h.logger.Errorw("Internal Server Error", "error", internalServerErr.Err, "message", internalServerErr.Message)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"type":  "internal_server",
		})

	default:
		// Log unexpected errors
		h.logger.Errorw("Unexpected Error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"type":  "unknown",
		})
	}
}

// HandleErrorFunc is a convenience function for error handling without creating an instance
// This requires passing a logger function
func HandleErrorFunc(c *gin.Context, err error, loggerFunc func(msg string, keysAndValues ...interface{})) {
	handler := &HTTPErrorHandler{
		logger: loggerAdapter{logFunc: loggerFunc},
	}
	handler.HandleError(c, err)
}

// loggerAdapter adapts a simple function to the Logger interface
type loggerAdapter struct {
	logFunc func(msg string, keysAndValues ...interface{})
}

func (l loggerAdapter) Errorw(msg string, keysAndValues ...interface{}) {
	l.logFunc(msg, keysAndValues...)
}
