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

