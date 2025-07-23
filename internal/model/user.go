// internal/model/user.go
package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string     `json:"name" validate:"required"`
	Email        string     `json:"email" gorm:"uniqueIndex" validate:"required,email"`
	Password     string     `json:"password" validate:"required,min=6"`
	Phone        string     `json:"phone" validate:"required"`
	RefreshToken string     `json:"-" gorm:"index"` // Store refresh token, exclude from JSON
	TokenExpiry  *time.Time `json:"-"`              // Track when refresh token expires
}
