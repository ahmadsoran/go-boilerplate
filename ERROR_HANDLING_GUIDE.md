# Enhanced Error Handling System

This document demonstrates the comprehensive error handling system implemented in the Go boilerplate project.

## Error Types Implemented

### 1. NotFoundError
**Usage:** When a resource is not found
```go
pkg.NewNotFoundError("User with ID %d not found", userID)
```
**HTTP Status:** 404 Not Found

### 2. InvalidInputError  
**Usage:** For general validation errors and malformed input
```go
pkg.NewInvalidInputError("Invalid user ID format")
```
**HTTP Status:** 400 Bad Request

### 3. DuplicateError
**Usage:** When trying to create a resource that already exists (unique constraint violations)
```go
pkg.NewDuplicateError("email", "User with email %s already exists", email)
```
**HTTP Status:** 409 Conflict
**Additional Fields:** `field` (the conflicting field)

### 4. ValidationError
**Usage:** For field-specific validation errors with context
```go
pkg.NewValidationError("password", "***", "Password must be at least 6 characters long")
```
**HTTP Status:** 400 Bad Request
**Additional Fields:** `field`, `value`

### 5. UnauthorizedError
**Usage:** For authentication failures
```go
pkg.NewUnauthorizedError("Invalid email or password")
```
**HTTP Status:** 401 Unauthorized

### 6. ForbiddenError
**Usage:** When user is authenticated but lacks permission
```go
pkg.NewForbiddenError("user_profile", "update", "You can only update your own profile")
```
**HTTP Status:** 403 Forbidden
**Additional Fields:** `resource`, `action`

### 7. ConflictError
**Usage:** For business logic conflicts
```go
pkg.NewConflictError("user_status", "Cannot delete user with active sessions")
```
**HTTP Status:** 409 Conflict
**Additional Fields:** `details`

### 8. RateLimitError
**Usage:** When rate limiting is applied
```go
pkg.NewRateLimitError(60, "Too many login attempts. Try again in %d seconds", 60)
```
**HTTP Status:** 429 Too Many Requests
**Headers:** `Retry-After` with retry time
**Additional Fields:** `retry_after`

### 9. ServiceUnavailableError
**Usage:** When external services are down
```go
pkg.NewServiceUnavailableError("email_service", err, "Email service is temporarily unavailable")
```
**HTTP Status:** 503 Service Unavailable
**Additional Fields:** `service`

### 10. TimeoutError
**Usage:** For request timeouts
```go
pkg.NewTimeoutError("database_query", 30, "Database query timed out after %d seconds", 30)
```
**HTTP Status:** 408 Request Timeout
**Additional Fields:** `operation`, `timeout`

### 11. PayloadTooLargeError
**Usage:** When request payload exceeds limits
```go
pkg.NewPayloadTooLargeError(1024*1024, 2*1024*1024, "File size %d exceeds maximum allowed size %d")
```
**HTTP Status:** 413 Payload Too Large
**Additional Fields:** `max_size`, `actual_size`

### 12. UnsupportedMediaTypeError
**Usage:** For unsupported content types
```go
supportedTypes := []string{"application/json", "application/xml"}
pkg.NewUnsupportedMediaTypeError("text/plain", supportedTypes, "Content type %s not supported")
```
**HTTP Status:** 415 Unsupported Media Type
**Additional Fields:** `received_type`, `supported_types`

### 13. BadGatewayError
**Usage:** When upstream services return errors
```go
pkg.NewBadGatewayError("payment_service", err, "Payment service returned an error")
```
**HTTP Status:** 502 Bad Gateway
**Additional Fields:** `upstream`

### 14. TooManyRequestsError
**Usage:** When external services rate limit us
```go
pkg.NewTooManyRequestsError("api_service", 120, "External API rate limit exceeded")
```
**HTTP Status:** 429 Too Many Requests
**Headers:** `Retry-After` with retry time
**Additional Fields:** `service`, `retry_after`

### 15. DatabaseConnectionError
**Usage:** For database connectivity issues
```go
pkg.NewDatabaseConnectionError("postgresql", err, "Failed to connect to database")
```
**HTTP Status:** 503 Service Unavailable
**Additional Fields:** `database`

### 16. MigrationError
**Usage:** For database migration failures
```go
pkg.NewMigrationError("001_create_users", err, "Migration failed")
```
**HTTP Status:** 500 Internal Server Error
**Additional Fields:** `migration`

### 17. ConfigurationError
**Usage:** For configuration-related errors
```go
pkg.NewConfigurationError("JWT_SECRET", "string", "JWT secret must be a non-empty string")
```
**HTTP Status:** 500 Internal Server Error
**Additional Fields:** `config_key`, `expected_type`

### 18. FileNotFoundError
**Usage:** For file system errors
```go
pkg.NewFileNotFoundError("/path/to/file", "Configuration file not found")
```
**HTTP Status:** 404 Not Found
**Additional Fields:** `file_path`

### 19. PermissionDeniedError
**Usage:** For file system permission errors
```go
pkg.NewPermissionDeniedError("/var/logs", "write", "Permission denied to write to log directory")
```
**HTTP Status:** 403 Forbidden
**Additional Fields:** `resource`, `operation`

### 20. NetworkError
**Usage:** For network connectivity issues
```go
pkg.NewNetworkError("api.example.com", 443, err, "Failed to connect to external API")
```
**HTTP Status:** 503 Service Unavailable
**Additional Fields:** `host`, `port`

### 21. CacheError
**Usage:** For caching system errors
```go
pkg.NewCacheError("redis", "user:123", err, "Failed to retrieve cached user data")
```
**HTTP Status:** 503 Service Unavailable
**Additional Fields:** `cache_type`, `key`

### 22. QueueError
**Usage:** For message queue errors
```go
pkg.NewQueueError("email_queue", "publish", err, "Failed to publish message to queue")
```
**HTTP Status:** 503 Service Unavailable
**Additional Fields:** `queue`, `operation`

### 23. ExternalAPIError
**Usage:** For third-party API errors
```go
pkg.NewExternalAPIError("stripe", "/v1/charges", 402, err, "Payment required")
```
**HTTP Status:** 502 Bad Gateway
**Additional Fields:** `api`, `endpoint`, `status_code`

### 24. InternalServerError
**Usage:** For unexpected internal errors
```go
pkg.NewInternalServerError(err, "Unexpected database error during user creation")
```
**HTTP Status:** 500 Internal Server Error

## Error Response Format

All errors now return structured JSON responses with consistent formatting:

```json
{
  "error": "User with email john@example.com already exists",
  "type": "duplicate",
  "field": "email"
}
```

### Common Response Fields:
- `error`: Human-readable error message
- `type`: Error type identifier for client-side handling
- Additional context fields specific to each error type

## Usage in Handlers

The `handleError` method in user handlers automatically maps error types to appropriate HTTP status codes and response formats:

```go
func (h *UserHandler) SomeMethod(c *gin.Context) {
    if err := h.svc.SomeOperation(ctx); err != nil {
        h.handleError(c, err) // Automatically handles all error types
        return
    }
    // Success response
}
```

## Security Considerations

1. **Information Disclosure Prevention:** 
   - `NotFoundError` for users is converted to `UnauthorizedError` in login to prevent username enumeration
   - Internal error details are logged but not exposed to clients

2. **Consistent Error Messages:**
   - Authentication errors use consistent messaging regardless of whether user exists

3. **Rate Limiting Headers:**
   - Proper `Retry-After` headers for rate limiting errors

## Logging

All errors are logged with appropriate context:
- Internal errors: Full error details logged
- Client errors: Request context logged
- Security errors: Additional security context logged

## Examples in Current Implementation

1. **Duplicate Email Registration:**
   ```go
   // In repository layer
   return pkg.NewDuplicateError("email", "User with email %s already exists", user.Email)
   ```

2. **Invalid Login:**
   ```go
   // In service layer  
   return pkg.NewUnauthorizedError("Invalid email or password")
   ```

3. **Expired Refresh Token:**
   ```go
   // In service layer
   return pkg.NewUnauthorizedError("Refresh token has expired")
   ```

4. **Validation Error:**
   ```go
   // In service layer
   return pkg.NewValidationError("password", "***", "Password must be at least 6 characters long")
   ```
