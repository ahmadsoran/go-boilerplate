# JWT Authentication Documentation

## Overview

This Go boilerplate now includes JWT (JSON Web Token) authentication with the following features:

- User registration and login
- Password hashing using bcrypt
- JWT token generation and validation
- Protected routes using middleware
- Context-based user information extraction

## Environment Variables

Add the following environment variables to your `.env` file:

```env
JWT_SECRET=your-secret-key-here-make-it-long-and-random
JWT_EXPIRY_HOURS=24
```

## Authentication Endpoints

### Register User
```bash
POST /api/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com", 
  "password": "password123",
  "phone": "+1234567890"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

### Login User
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe", 
    "email": "john@example.com",
    "phone": "+1234567890",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

## Protected Routes

The following routes require authentication via JWT token:

- `GET /api/users/:id` - Get user by ID
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

### Using Protected Routes

Include the JWT token in the Authorization header:

```bash
GET /api/users/1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Code Structure

### JWT Manager (`internal/pkg/jwt.go`)
- `GenerateToken(userID, email)` - Creates JWT tokens
- `ValidateToken(tokenString)` - Validates and parses JWT tokens

### Password Utilities (`internal/pkg/password.go`)
- `HashPassword(password)` - Hashes passwords using bcrypt
- `CheckPassword(hashedPassword, password)` - Verifies passwords

### Auth Middleware (`internal/middleware/auth.middleware.go`)
- `AuthMiddleware(jwtManager)` - Validates JWT tokens
- `GetUserIDFromContext(c)` - Extracts user ID from Gin context
- `GetUserEmailFromContext(c)` - Extracts user email from Gin context

### Service Layer
- `RegisterUser()` - Handles user registration with password hashing
- `LoginUser()` - Authenticates users and validates passwords
- `GetUserByEmail()` - Retrieves users by email

## Testing the API

1. Start the server:
```bash
go run cmd/main.go
```

2. Register a new user:
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123", 
    "phone": "+1234567890"
  }'
```

3. Login with the user:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

4. Use the token to access protected routes:
```bash
curl -X GET http://localhost:8080/api/users/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Security Considerations

- JWT secret should be a long, random string
- Store JWT secret as an environment variable, never in code
- Tokens expire after the configured time (default: 24 hours)
- Passwords are hashed using bcrypt with default cost
- Consider implementing token refresh for better UX
- Add rate limiting for authentication endpoints in production
