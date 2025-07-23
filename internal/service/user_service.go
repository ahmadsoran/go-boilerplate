// internal/service/user_service.go
package service

import (
	"context"

	"your_project/internal/logger"
	userlogger "your_project/internal/logger/user-logger"
	"your_project/internal/model"
	"your_project/internal/repository"

	"gorm.io/gorm"
)

type UserService interface {
	GetUser(ctx context.Context, id uint) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id uint) error
}

type userService struct {
	repo repository.UserRepository
	db   *gorm.DB
}

func NewUserService(repo repository.UserRepository, db *gorm.DB) UserService {
	return &userService{repo, db}
}

func (s *userService) GetUser(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// Propagate repository errors, which are already custom errors
		return nil, err
	}
	return user, nil
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) error {
	logger := logger.APILog
	logger.Info("Creating user", "user", user.Name)
	// Add any business logic validation here and return pkg.NewInvalidInputError if needed

	if err := s.repo.Create(ctx, user); err != nil {
		// Propagate repository errors
		return err
	}
	return nil
}

func (s *userService) UpdateUser(ctx context.Context, user *model.User) error {
	// Add any business logic validation here and return pkg.NewInvalidInputError if needed
	logger := userlogger.GetUserLogger(user.ID)

	// Pass the context to the transaction
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repoTx := s.repo.WithTx(tx)
		// Pass the context to repository calls within the transaction
		oldUser, err := repoTx.GetByID(ctx, user.ID)
		if err != nil {
			// Propagate repository errors (e.g., NotFoundError)
			logger.Errorw("User not found during update transaction", "userID", user.ID, "error", err)
			return err
		}
		oldUser.Email = user.Email
		oldUser.Name = user.Name
		// Pass the context to repository calls within the transaction
		if err := repoTx.Update(ctx, oldUser); err != nil {
			// Propagate repository errors
			logger.Errorw("Update failed during transaction", "userID", oldUser.ID, "error", err)
			return err
		}
		logger.Info("User updated", "user", oldUser.ID)
		return nil
	})
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	logger := userlogger.GetUserLogger(id)
	logger.Info("Deleting user", "userID", id)
	// Add any business logic validation here and return pkg.NewInvalidInputError if needed

	if err := s.repo.Delete(ctx, id); err != nil {
		// Propagate repository errors
		return err
	}
	return nil
}
