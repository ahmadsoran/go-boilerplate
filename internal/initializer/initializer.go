// internal/initializer/initializer.go
package initializer

import (
	"gorm.io/gorm"
	"your_project/internal/api"
	"your_project/internal/repository"
	"your_project/internal/service"
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
	User   *api.UserHandler
	Health *api.HealthHandler
	// Add other handlers here
}

func NewHandlerContainer(svcs *ServiceContainer, db *gorm.DB) *HandlerContainer {
	return &HandlerContainer{
		User:   api.NewUserHandler(svcs.User),
		Health: api.NewHealthHandler(db),
		// Add other handlers here
	}
}
