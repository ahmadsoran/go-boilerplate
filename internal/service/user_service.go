// internal/service/user_service.go
package service

import (
	"context"
	"time"

	"your_project/internal/logger"
	userlogger "your_project/internal/logger/user-logger"
	"your_project/internal/model"
	"your_project/internal/pkg"
	"your_project/internal/repository"

	"gorm.io/gorm"
)

type UserService interface {
	GetUser(ctx context.Context, id uint) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id uint) error
	RegisterUser(ctx context.Context, user *model.User) error
	LoginUser(ctx context.Context, email, password string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.User, error)
	UpdateRefreshToken(ctx context.Context, userID uint, refreshToken string, expiry time.Time) error
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

// RegisterUser creates a new user with hashed password
func (s *userService) RegisterUser(ctx context.Context, user *model.User) error {
	logger := logger.APILog
	logger.Info("Registering user", "email", user.Email)

	// Validate email format (additional business logic validation)
	if user.Email == "" {
		return pkg.NewValidationError("email", user.Email, "Email is required")
	}

	// Validate password strength
	if len(user.Password) < 6 {
		return pkg.NewValidationError("password", "***", "Password must be at least 6 characters long")
	}

	// Validate name
	if user.Name == "" {
		return pkg.NewValidationError("name", user.Name, "Name is required")
	}

	// Validate phone
	if user.Phone == "" {
		return pkg.NewValidationError("phone", user.Phone, "Phone number is required")
	}

	// Hash the password before storing
	hashedPassword, err := pkg.HashPassword(user.Password)
	if err != nil {
		return pkg.NewInternalServerError(err, "Failed to hash password")
	}
	user.Password = hashedPassword

	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}
	return nil
}

// LoginUser authenticates a user with email and password and returns token pair
func (s *userService) LoginUser(ctx context.Context, email, password string) (*model.User, error) {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		// Convert NotFoundError to UnauthorizedError to avoid revealing user existence
		return nil, pkg.NewUnauthorizedError("Invalid email or password")
	}

	// Check password
	if err := pkg.CheckPassword(user.Password, password); err != nil {
		return nil, pkg.NewUnauthorizedError("Invalid email or password")
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// RefreshToken validates a refresh token and returns the associated user
func (s *userService) RefreshToken(ctx context.Context, refreshToken string) (*model.User, error) {
	user, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		// Convert NotFoundError to UnauthorizedError for security
		return nil, pkg.NewUnauthorizedError("Invalid or expired refresh token")
	}

	// Check if the refresh token has expired
	if user.TokenExpiry != nil && time.Now().After(*user.TokenExpiry) {
		return nil, pkg.NewUnauthorizedError("Refresh token has expired")
	}

	return user, nil
}

// UpdateRefreshToken updates a user's refresh token and expiry
func (s *userService) UpdateRefreshToken(ctx context.Context, userID uint, refreshToken string, expiry time.Time) error {
	logger := userlogger.GetUserLogger(userID)
	logger.Info("Updating refresh token for user", "userID", userID)

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repoTx := s.repo.WithTx(tx)

		user, err := repoTx.GetByID(ctx, userID)
		if err != nil {
			logger.Errorw("User not found during refresh token update", "userID", userID, "error", err)
			return err
		}

		user.RefreshToken = refreshToken
		user.TokenExpiry = &expiry

		if err := repoTx.Update(ctx, user); err != nil {
			logger.Errorw("Failed to update refresh token", "userID", userID, "error", err)
			return err
		}

		logger.Info("Refresh token updated successfully", "userID", userID)
		return nil
	})
}
