// internal/initializer/initializer.go
package initializer

import (
	"your_project/internal/api/handlers"
	"your_project/internal/repository"
	"your_project/internal/service"

	"gorm.io/gorm"
)

type RepositoryContainer struct {
	User repository.UserRepository
	// Add other repositories here
}

func NewRepositoryContainer(db *gorm.DB) *RepositoryContainer {
	return &RepositoryContainer{
		User: repository.NewUserRepository(db),
		// Add other repositories here
	}
}

type ServiceContainer struct {
	User service.UserService
	// Add other services here
}

func NewServiceContainer(repos *RepositoryContainer, db *gorm.DB) *ServiceContainer {
	return &ServiceContainer{
		User: service.NewUserService(repos.User, db),
		// Add other services here
	}
}

type HandlerContainer struct {
	User   *handlers.UserHandler
	Health *handlers.HealthHandler
	// Add other handlers here
}

func NewHandlerContainer(svcs *ServiceContainer, db *gorm.DB) *HandlerContainer {
	return &HandlerContainer{
		User:   handlers.NewUserHandler(svcs.User),
		Health: handlers.NewHealthHandler(db),
		// Add other handlers here
	}
}
