// internal/repository/user_repo.go
package repository

import (
	"context"
	"errors"
	"strings"

	"your_project/internal/model"
	"your_project/internal/pkg"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
	WithTx(tx *gorm.DB) UserRepository
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) WithTx(tx *gorm.DB) UserRepository {
	return &userRepository{tx}
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User

	// Pass the context to the GORM query
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NewNotFoundError("user with ID %d not found", id)
		}
		return nil, pkg.NewInternalServerError(err, "failed to get user by ID %d", id)
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	// Pass the context to the GORM query
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NewNotFoundError("user with email %s not found", email)
		}
		return nil, pkg.NewInternalServerError(err, "failed to get user by email %s", email)
	}

	return &user, nil
}

func (r *userRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*model.User, error) {
	var user model.User

	// Pass the context to the GORM query
	if err := r.db.WithContext(ctx).Where("refresh_token = ?", refreshToken).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NewNotFoundError("user with refresh token not found")
		}
		return nil, pkg.NewInternalServerError(err, "failed to get user by refresh token")
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	// Pass the context to the GORM query
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		// Check for duplicate key violations (unique constraints)
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") ||
			strings.Contains(err.Error(), "Duplicate entry") {
			// Determine which field caused the duplicate error
			if strings.Contains(err.Error(), "email") || strings.Contains(err.Error(), "users_email_key") {
				return pkg.NewDuplicateError("email", "User with email %s already exists", user.Email)
			}
			return pkg.NewDuplicateError("unknown", "Duplicate entry detected")
		}
		return pkg.NewInternalServerError(err, "failed to create user")
	}
	return nil
}
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	// Pass the context to the GORM query
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		// Check for duplicate key violations (unique constraints)
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") ||
			strings.Contains(err.Error(), "Duplicate entry") {
			// Determine which field caused the duplicate error
			if strings.Contains(err.Error(), "email") || strings.Contains(err.Error(), "users_email_key") {
				return pkg.NewDuplicateError("email", "Email %s is already taken by another user", user.Email)
			}
			return pkg.NewDuplicateError("unknown", "Duplicate entry detected during update")
		}
		return pkg.NewInternalServerError(err, "failed to update user with ID %d", user.ID)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	// Pass the context to the GORM query
	if err := r.db.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {
		return pkg.NewInternalServerError(err, "failed to delete user with ID %d", id)
	}
	return nil
}
