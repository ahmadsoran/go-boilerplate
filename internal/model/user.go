// internal/model/user.go
package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" gorm:"uniqueIndex" validate:"required,email"`
	Password string `json:"-" validate:"required,min=6"`
	Phone    string `json:"-" validate:"required,min=14"`
}
