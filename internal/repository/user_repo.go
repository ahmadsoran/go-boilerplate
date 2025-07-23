// internal/repository/user_repo.go
package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"your_project/internal/model"
	"your_project/internal/pkg"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uint) (*model.User, error)
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

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	// Pass the context to the GORM query
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return pkg.NewInternalServerError(err, "failed to create user")
	}
	return nil
}
unc (r *userRepository) Update(ctx context.Context, user *model.User) error {
	// Pass the context to the GORM query
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		// Consider checking for specific errors like unique constraint violations
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
