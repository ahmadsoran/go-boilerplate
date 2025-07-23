📁 Folder Structure
plaintext
Copy
Edit
your_project/
│
├── cmd/                  # Main application entry point
│   └── main.go
│
├── configs/              # Configuration files (env, yaml, json, etc.)
│   └── config.go
│
├── internal/
│   ├── api/              # HTTP handlers / controllers
│   │   └── user_handler.go
│   │
│   ├── service/          # Business logic layer
│   │   └── user_service.go
│   │
│   ├── repository/       # Data access layer
│   │   └── user_repo.go
│   │
│   ├── model/            # GORM models and DTOs
│   │   └── user.go
│   │
│   ├── db/               # DB initialization & migrations
│   │   └── db.go
│   │
│   └── logger/           # Logging setup
│       └── logger.go
│
├── go.mod
├── go.sum
└── README.md
📝 Example: User CRUD with Transaction & Logging
1. User Model (internal/model/user.go)
go
Copy
Edit
package model

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Name     string `json:"name"`
    Email    string `json:"email" gorm:"uniqueIndex"`
    Password string `json:"-"`
}
2. Repository Layer (internal/repository/user_repo.go)
go
Copy
Edit
package repository

import (
    "your_project/internal/model"
    "gorm.io/gorm"
)

type UserRepository interface {
    GetByID(id uint) (*model.User, error)
    Create(user *model.User) error
    Update(user *model.User) error
    Delete(id uint) error
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

func (r *userRepository) GetByID(id uint) (*model.User, error) {
    var user model.User
    if err := r.db.First(&user, id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) Update(user *model.User) error {
    return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
    return r.db.Delete(&model.User{}, id).Error
}
3. Service Layer (internal/service/user_service.go)
go
Copy
Edit
package service

import (
    "your_project/internal/model"
    "your_project/internal/repository"
    "your_project/internal/logger"
    "gorm.io/gorm"
)

type UserService interface {
    GetUser(id uint) (*model.User, error)
    CreateUser(user *model.User) error
    UpdateUser(user *model.User) error
    DeleteUser(id uint) error
}

type userService struct {
    repo repository.UserRepository
    db   *gorm.DB
}

func NewUserService(repo repository.UserRepository, db *gorm.DB) UserService {
    return &userService{repo, db}
}

func (s *userService) GetUser(id uint) (*model.User, error) {
    return s.repo.GetByID(id)
}

func (s *userService) CreateUser(user *model.User) error {
    logger.Log.Info("Creating user", "user", user.Name)
    return s.repo.Create(user)
}

// Transaction example: update user email safely
func (s *userService) UpdateUser(user *model.User) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        repoTx := s.repo.WithTx(tx)
        oldUser, err := repoTx.GetByID(user.ID)
        if err != nil {
            logger.Log.Error("User not found", "error", err)
            return err
        }
        oldUser.Email = user.Email
        oldUser.Name = user.Name
        if err := repoTx.Update(oldUser); err != nil {
            logger.Log.Error("Update failed", "error", err)
            return err
        }
        logger.Log.Info("User updated", "user", oldUser.ID)
        return nil
    })
}

func (s *userService) DeleteUser(id uint) error {
    logger.Log.Info("Deleting user", "userID", id)
    return s.repo.Delete(id)
}
4. API Layer / Handler (internal/api/user_handler.go)
go
Copy
Edit
package api

import (
    "net/http"
    "strconv"
    "your_project/internal/model"
    "your_project/internal/service"
    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
    return &UserHandler{svc}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
    users := r.Group("/users")
    {
        users.POST("", h.CreateUser)
        users.GET("/:id", h.GetUser)
        users.PUT("/:id", h.UpdateUser)
        users.DELETE("/:id", h.DeleteUser)
    }
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var user model.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.svc.CreateUser(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    user, err := h.svc.GetUser(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var user model.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    user.ID = uint(id)
    if err := h.svc.UpdateUser(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if err := h.svc.DeleteUser(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusNoContent, nil)
}
5. Logger Setup (internal/logger/logger.go)
go
Copy
Edit
package logger

import (
    "log"
    "go.uber.org/zap"
)

var Log *zap.SugaredLogger

func Init() {
    logger, err := zap.NewDevelopment()
    if err != nil {
        log.Fatalf("Can't init zap logger: %v", err)
    }
    Log = logger.Sugar()
}
6. DB Setup (internal/db/db.go)
go
Copy
Edit
package db

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "your_project/internal/model"
)

func Init() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    db.AutoMigrate(&model.User{})
    return db, nil
}
7. Main Entrypoint (cmd/main.go)
go
Copy
Edit
package main

import (
    "github.com/gin-gonic/gin"
    "your_project/internal/api"
    "your_project/internal/service"
    "your_project/internal/repository"
    "your_project/internal/db"
    "your_project/internal/logger"
)

func main() {
    logger.Init()
    dbConn, err := db.Init()
    if err != nil {
        logger.Log.Fatalw("DB connection failed", "error", err)
    }

    userRepo := repository.NewUserRepository(dbConn)
    userService := service.NewUserService(userRepo, dbConn)
    userHandler := api.NewUserHandler(userService)

    r := gin.Default()
    userHandler.RegisterRoutes(r)

    logger.Log.Info("Server started at :8080")
    r.Run(":8080")
}
📄 MDF/README.md (Documentation Sample)
markdown
Copy
Edit
# Go Gin-GORM Clean Architecture Example

## Folder Structure

- `cmd/` - Entry point for the application (`main.go`)
- `configs/` - App configs
- `internal/`
    - `api/` - HTTP handlers/controllers (routes, request/response)
    - `service/` - Business logic and transaction control
    - `repository/` - Database access via GORM, all raw DB queries go here
    - `model/` - GORM models (structs mapped to tables)
    - `db/` - DB setup, migrations
    - `logger/` - Logging setup, using Uber Zap

## Key Features

- **Separation of Concerns:** Each layer has one responsibility; no leaking business logic to controller or repository.
- **Repository Pattern:** All data access abstracted; easy to swap DB or mock in tests.
- **Service Layer:** All business logic, including transactions, validation, etc.
- **Transaction Example:** Service wraps DB operations in a GORM transaction for safe updates.
- **Logging:** All actions, errors, and business events logged with context using Uber Zap.
- **CRUD Endpoints:** Standard RESTful routes with JSON request/response.
- **Error Handling:** Consistent HTTP error responses.

## How to Run

```bash
go mod tidy
go run ./cmd/main.go
Example API Endpoints
POST /users - Create user

GET /users/:id - Get user

PUT /users/:id - Update user (with transaction, logs)

DELETE /users/:id - Delete user

Customization
Swap SQLite for PostgreSQL/MySQL in internal/db/db.go.

Add more services/repositories/models as needed.

Plug your configs/environment variables in configs/.

