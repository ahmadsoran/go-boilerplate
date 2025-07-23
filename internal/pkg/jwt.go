// internal/pkg/jwt.go
package pkg

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the JWT claims
type JWTClaims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// JWTManager handles JWT operations
type JWTManager struct {
	secretKey          string
	expiryHours        int
	refreshExpiryHours int
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey string, expiryHours int) *JWTManager {
	return &JWTManager{
		secretKey:          secretKey,
		expiryHours:        expiryHours,
		refreshExpiryHours: expiryHours * 7, // Refresh tokens last 7 times longer
	}
}

// TokenPair represents both access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GenerateTokenPair generates both access and refresh tokens for a user
func (j *JWTManager) GenerateTokenPair(userID uint, email string) (*TokenPair, error) {
	accessToken, err := j.generateToken(userID, email, "access", time.Duration(j.expiryHours)*time.Hour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.generateToken(userID, email, "refresh", time.Duration(j.refreshExpiryHours)*time.Hour)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GenerateToken generates a new JWT token for a user (backwards compatibility)
func (j *JWTManager) GenerateToken(userID uint, email string) (string, error) {
	return j.generateToken(userID, email, "access", time.Duration(j.expiryHours)*time.Hour)
}

// generateToken is the internal method for generating tokens
func (j *JWTManager) generateToken(userID uint, email, tokenType string, expiry time.Duration) (string, error) {
	claims := JWTClaims{
		UserID:    userID,
		Email:     email,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// RefreshExpiryHours returns the refresh token expiry hours
func (j *JWTManager) RefreshExpiryHours() int {
	return j.refreshExpiryHours
}
