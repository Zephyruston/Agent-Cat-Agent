# 🐹 Go Large System Development Rules

> **A comprehensive rule system for building production-ready, scalable Go applications with consistent architecture patterns, proper error handling, and modern concurrency patterns.**

## 📚 Rule System Overview

This rule system automatically selects appropriate architectural patterns, dependency management strategies, and coding standards based on project complexity and functional requirements. It ensures maintainability and scalability of large Go systems through:

- **Complexity-driven architecture decisions**
- **Structured error handling patterns** (pkg/errors for libs, fmt.Errorf for simple cases)
- **Frontend-friendly serialization** (CamelCase JSON output with proper tags)
- **Builder patterns** for complex constructors
- **Modern concurrency patterns** using channels, sync primitives, and worker pools
- **Structured module management** with proper workspace organization

## 🏗️ Rule Architecture

```
.cursor/go/
├── core/                           # Core rule system
│   ├── main.mdc                   # Main rule entry point
│   ├── complexity-detection.mdc   # Project complexity analysis
│   └── module-management.mdc      # Module organization
├── modules/                        # Feature-specific rules
│   ├── error-handling/            # Error handling patterns
│   │   ├── structured-errors.mdc  # Custom error types
│   │   └── error-wrapping.mdc     # Error context patterns
│   ├── serialization/             # JSON/serialization rules
│   │   └── json-patterns.mdc      # JSON tag conventions
│   ├── builder-patterns/          # Constructor patterns
│   │   └── functional-options.mdc # Functional options pattern
│   ├── concurrency/               # Concurrency patterns
│   │   └── modern-patterns.mdc    # Channels, workers, sync
│   ├── database/                  # Database patterns
│   │   └── repository-patterns.mdc # Repository with interfaces
│   ├── testing/                   # Testing patterns
│   │   └── test-organization.mdc  # Test structure and patterns
│   └── [other modules...]         # Additional patterns
└── workflows/                      # Development workflows
    └── development-workflow.mdc    # Complete dev lifecycle
```

## 🎯 Quick Start

### 1. Project Initialization

The rule system automatically detects project complexity:

- **Simple (0-6 points)**: Single module, basic patterns
- **Medium (7-15 points)**: Multi-package, structured errors, builders
- **Complex (16+ points)**: Multi-module, advanced patterns, full concurrency

### 2. Automatic Rule Loading

Rules are automatically loaded based on:
- Project complexity score
- Detected dependencies (gin, gorm, etc.)
- Package types (cmd vs pkg vs internal)
- Feature requirements

### 3. Example Project Structure

For a complex web application:

```
my-app/
├── go.mod                  # Main module
├── go.sum                  # Dependency checksums
├── cmd/                    # Application entry points
│   ├── api-server/
│   │   └── main.go
│   └── cli-tool/
│       └── main.go
├── internal/               # Private application code
│   ├── app/               # Application logic
│   ├── domain/            # Business domain
│   └── infrastructure/    # External dependencies
├── pkg/                   # Public library code
│   ├── client/            # API clients
│   └── models/            # Shared models
├── api/                   # API definitions
│   └── openapi/
├── web/                   # Web assets (if applicable)
├── configs/               # Configuration files
├── scripts/               # Build and deployment scripts
├── deployments/           # Deployment configurations
└── docs/                  # Documentation
```

## 📋 Rule Categories

### 🏗️ Core Architecture Rules

**Module Management** (`core/module-management.mdc`)
- Use Go modules for dependency management
- Proper module naming conventions
- Semantic versioning for releases

**Complexity Detection** (`core/complexity-detection.mdc`)
- Automated project complexity assessment
- Rule selection based on complexity score
- Adaptive architectural recommendations

### 🚨 Error Handling Rules

**Structured Errors** (`modules/error-handling/structured-errors.mdc`)
- Custom error types implementing error interface
- Error codes and categories
- Proper error context preservation

**Error Wrapping** (`modules/error-handling/error-wrapping.mdc`)
- Use `fmt.Errorf` with `%w` verb for wrapping
- Structured error context with additional data
- Error chain inspection with `errors.Is` and `errors.As`

### 📦 Serialization Rules

**JSON Patterns** (`modules/serialization/json-patterns.mdc`)
- All structs must use proper JSON tags
- CamelCase for external APIs
- Consistent field naming conventions

### 🏗️ Constructor Rules

**Functional Options** (`modules/builder-patterns/functional-options.mdc`)
- Use functional options for complex constructors
- Proper default values and validation
- Type-safe construction patterns

### 🏗️ Domain Organization Rules

**Domain-Driven Structure** (`modules/domain-organization/domain-driven-structure.mdc`)
- Organize by business domain in internal/domain/
- Use meaningful package names
- Proper separation of concerns

### 🗄️ Database Rules

**Repository Patterns** (`modules/database/repository-patterns.mdc`)
- Interface-based repository pattern
- Proper transaction handling
- Database error handling

### 🧪 Testing Rules

**Test Organization** (`modules/testing/test-organization.mdc`)
- Table-driven tests for multiple cases
- Proper test structure and naming
- Integration test patterns

### ⚡ Concurrency Rules

**Modern Patterns** (`modules/concurrency/modern-patterns.mdc`)
- Channel-based communication
- Worker pool patterns
- Proper context usage for cancellation

## 🔧 Implementation Examples

### Error Handling

```go
// pkg/errors/types.go - Custom error types
package errors

import (
    "fmt"
    "net/http"
)

// ErrorCode represents different types of errors
type ErrorCode string

const (
    CodeValidation   ErrorCode = "VALIDATION_ERROR"
    CodeNotFound     ErrorCode = "NOT_FOUND"
    CodeUnauthorized ErrorCode = "UNAUTHORIZED"
    CodeInternal     ErrorCode = "INTERNAL_ERROR"
)

// AppError represents application-specific errors
type AppError struct {
    Code    ErrorCode `json:"code"`
    Message string    `json:"message"`
    Details any       `json:"details,omitempty"`
    Cause   error     `json:"-"`
}

func (e *AppError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
    return e.Cause
}

func (e *AppError) HTTPStatus() int {
    switch e.Code {
    case CodeValidation:
        return http.StatusBadRequest
    case CodeNotFound:
        return http.StatusNotFound
    case CodeUnauthorized:
        return http.StatusUnauthorized
    default:
        return http.StatusInternalServerError
    }
}

// Error constructors
func NewValidationError(message string, details any) *AppError {
    return &AppError{
        Code:    CodeValidation,
        Message: message,
        Details: details,
    }
}

func NewNotFoundError(resource string) *AppError {
    return &AppError{
        Code:    CodeNotFound,
        Message: fmt.Sprintf("%s not found", resource),
    }
}

func WrapError(err error, code ErrorCode, message string) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
        Cause:   err,
    }
}
```

### Serialization

```go
// pkg/models/user.go - JSON serialization patterns
package models

import (
    "time"
)

// User represents user data with proper JSON tags
type User struct {
    ID        string    `json:"id" db:"id"`
    FirstName string    `json:"firstName" db:"first_name"`
    LastName  string    `json:"lastName" db:"last_name"`
    Email     string    `json:"email" db:"email"`
    CreatedAt time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
    
    // Sensitive fields excluded from JSON
    Password string `json:"-" db:"password"`
}

// UserCreateRequest represents user creation request
type UserCreateRequest struct {
    FirstName string `json:"firstName" validate:"required,min=2,max=50"`
    LastName  string `json:"lastName" validate:"required,min=2,max=50"`
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8"`
}

// UserResponse represents user response (without sensitive data)
type UserResponse struct {
    ID        string    `json:"id"`
    FirstName string    `json:"firstName"`
    LastName  string    `json:"lastName"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"createdAt"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
    return &UserResponse{
        ID:        u.ID,
        FirstName: u.FirstName,
        LastName:  u.LastName,
        Email:     u.Email,
        CreatedAt: u.CreatedAt,
    }
}
```

### Functional Options Pattern

```go
// internal/infrastructure/database/config.go - Database configuration
package database

import (
    "time"
)

// Config holds database configuration
type Config struct {
    Host            string
    Port            int
    Username        string
    Password        string
    Database        string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    SSLMode         string
    TimeZone        string
}

// Option defines functional option type
type Option func(*Config)

// NewConfig creates new database config with options
func NewConfig(host, username, password, database string, opts ...Option) *Config {
    // Set defaults
    config := &Config{
        Host:            host,
        Port:            5432,
        Username:        username,
        Password:        password,
        Database:        database,
        MaxOpenConns:    25,
        MaxIdleConns:    25,
        ConnMaxLifetime: 5 * time.Minute,
        SSLMode:         "require",
        TimeZone:        "UTC",
    }
    
    // Apply options
    for _, opt := range opts {
        opt(config)
    }
    
    return config
}

// Functional options
func WithPort(port int) Option {
    return func(c *Config) {
        c.Port = port
    }
}

func WithMaxConnections(max int) Option {
    return func(c *Config) {
        c.MaxOpenConns = max
        c.MaxIdleConns = max
    }
}

func WithSSLMode(mode string) Option {
    return func(c *Config) {
        c.SSLMode = mode
    }
}

func WithTimeZone(tz string) Option {
    return func(c *Config) {
        c.TimeZone = tz
    }
}

// Usage example:
// config := NewConfig("localhost", "user", "pass", "mydb",
//     WithPort(5433),
//     WithMaxConnections(50),
//     WithSSLMode("disable"),
// )
```

### Domain Organization

```go
// internal/domain/user/user.go - Domain entity
package user

import (
    "time"
    "github.com/google/uuid"
)

// User represents the user domain entity
type User struct {
    id        uuid.UUID
    firstName string
    lastName  string
    email     string
    createdAt time.Time
    updatedAt time.Time
}

// NewUser creates a new user with validation
func NewUser(firstName, lastName, email string) (*User, error) {
    if firstName == "" {
        return nil, NewValidationError("firstName is required")
    }
    if lastName == "" {
        return nil, NewValidationError("lastName is required")
    }
    if !isValidEmail(email) {
        return nil, NewValidationError("invalid email format")
    }
    
    now := time.Now()
    return &User{
        id:        uuid.New(),
        firstName: firstName,
        lastName:  lastName,
        email:     email,
        createdAt: now,
        updatedAt: now,
    }, nil
}

// Getters
func (u *User) ID() uuid.UUID { return u.id }
func (u *User) FirstName() string { return u.firstName }
func (u *User) LastName() string { return u.lastName }
func (u *User) Email() string { return u.email }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

// Business methods
func (u *User) UpdateName(firstName, lastName string) error {
    if firstName == "" {
        return NewValidationError("firstName is required")
    }
    if lastName == "" {
        return NewValidationError("lastName is required")
    }
    
    u.firstName = firstName
    u.lastName = lastName
    u.updatedAt = time.Now()
    return nil
}

func (u *User) FullName() string {
    return u.firstName + " " + u.lastName
}

// internal/domain/user/repository.go - Repository interface
package user

import (
    "context"
    "github.com/google/uuid"
)

// Repository defines the user repository interface
type Repository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id uuid.UUID) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, limit, offset int) ([]*User, error)
}

// internal/domain/user/service.go - Domain service
package user

import (
    "context"
    "fmt"
    "github.com/google/uuid"
)

// Service provides user business logic
type Service struct {
    repo Repository
}

// NewService creates a new user service
func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

// CreateUser creates a new user
func (s *Service) CreateUser(ctx context.Context, firstName, lastName, email string) (*User, error) {
    // Check if user already exists
    existing, err := s.repo.GetByEmail(ctx, email)
    if err == nil && existing != nil {
        return nil, NewValidationError("user with this email already exists")
    }
    
    // Create new user
    user, err := NewUser(firstName, lastName, email)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    // Save to repository
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to save user: %w", err)
    }
    
    return user, nil
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return user, nil
}
```

### Database Repository Pattern

```go
// internal/infrastructure/database/user_repository.go
package database

import (
    "context"
    "database/sql"
    "fmt"
    
    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    
    "myapp/internal/domain/user"
)

// userRepository implements user.Repository interface
type userRepository struct {
    db *sqlx.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sqlx.DB) user.Repository {
    return &userRepository{db: db}
}

// userEntity represents user data in database
type userEntity struct {
    ID        uuid.UUID `db:"id"`
    FirstName string    `db:"first_name"`
    LastName  string    `db:"last_name"`
    Email     string    `db:"email"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

// toDomain converts database entity to domain entity
func (e *userEntity) toDomain() *user.User {
    // Use reflection or constructor to recreate domain entity
    // This is a simplified example
    u, _ := user.NewUser(e.FirstName, e.LastName, e.Email)
    return u
}

// fromDomain converts domain entity to database entity
func (r *userRepository) fromDomain(u *user.User) *userEntity {
    return &userEntity{
        ID:        u.ID(),
        FirstName: u.FirstName(),
        LastName:  u.LastName(),
        Email:     u.Email(),
        CreatedAt: u.CreatedAt(),
        UpdatedAt: u.UpdatedAt(),
    }
}

// Create implements user.Repository
func (r *userRepository) Create(ctx context.Context, u *user.User) error {
    entity := r.fromDomain(u)
    
    query := `
        INSERT INTO users (id, first_name, last_name, email, created_at, updated_at)
        VALUES (:id, :first_name, :last_name, :email, :created_at, :updated_at)
    `
    
    _, err := r.db.NamedExecContext(ctx, query, entity)
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}

// GetByID implements user.Repository
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
    var entity userEntity
    
    query := `SELECT * FROM users WHERE id = $1`
    err := r.db.GetContext(ctx, &entity, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, user.NewNotFoundError("user")
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return entity.toDomain(), nil
}

// GetByEmail implements user.Repository
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
    var entity userEntity
    
    query := `SELECT * FROM users WHERE email = $1`
    err := r.db.GetContext(ctx, &entity, query, email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, user.NewNotFoundError("user")
        }
        return nil, fmt.Errorf("failed to get user by email: %w", err)
    }
    
    return entity.toDomain(), nil
}
```

### Testing Patterns

```go
// internal/domain/user/service_test.go - Table-driven tests
package user_test

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    
    "myapp/internal/domain/user"
)

// MockRepository is a mock implementation of user.Repository
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, u *user.User) error {
    args := m.Called(ctx, u)
    return args.Error(0)
}

func (m *MockRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
    args := m.Called(ctx, email)
    return args.Get(0).(*user.User), args.Error(1)
}

// Add other methods...

func TestService_CreateUser(t *testing.T) {
    ctx := context.Background()
    
    tests := []struct {
        name      string
        firstName string
        lastName  string
        email     string
        setupMock func(*MockRepository)
        wantErr   bool
        errType   error
    }{
        {
            name:      "successful user creation",
            firstName: "John",
            lastName:  "Doe",
            email:     "john@example.com",
            setupMock: func(mr *MockRepository) {
                mr.On("GetByEmail", ctx, "john@example.com").Return((*user.User)(nil), user.NewNotFoundError("user"))
                mr.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(nil)
            },
            wantErr: false,
        },
        {
            name:      "user already exists",
            firstName: "Jane",
            lastName:  "Doe",
            email:     "jane@example.com",
            setupMock: func(mr *MockRepository) {
                existingUser, _ := user.NewUser("Jane", "Doe", "jane@example.com")
                mr.On("GetByEmail", ctx, "jane@example.com").Return(existingUser, nil)
            },
            wantErr: true,
            errType: &user.ValidationError{},
        },
        {
            name:      "invalid email",
            firstName: "Bob",
            lastName:  "Smith",
            email:     "invalid-email",
            setupMock: func(mr *MockRepository) {
                // No mock setup needed as validation happens before repository call
            },
            wantErr: true,
            errType: &user.ValidationError{},
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockRepo := new(MockRepository)
            tt.setupMock(mockRepo)
            service := user.NewService(mockRepo)
            
            // Act
            result, err := service.CreateUser(ctx, tt.firstName, tt.lastName, tt.email)
            
            // Assert
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, result)
                if tt.errType != nil {
                    assert.IsType(t, tt.errType, err)
                }
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, result)
                assert.Equal(t, tt.firstName, result.FirstName())
                assert.Equal(t, tt.lastName, result.LastName())
                assert.Equal(t, tt.email, result.Email())
            }
            
            // Verify mock expectations
            mockRepo.AssertExpectations(t)
        })
    }
}

// Benchmark test example
func BenchmarkService_CreateUser(b *testing.B) {
    ctx := context.Background()
    mockRepo := new(MockRepository)
    mockRepo.On("GetByEmail", ctx, mock.AnythingOfType("string")).Return((*user.User)(nil), user.NewNotFoundError("user"))
    mockRepo.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(nil)
    
    service := user.NewService(mockRepo)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        email := fmt.Sprintf("user%d@example.com", i)
        _, err := service.CreateUser(ctx, "John", "Doe", email)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### Concurrency Patterns

```go
// internal/infrastructure/worker/pool.go - Worker pool pattern
package worker

import (
    "context"
    "fmt"
    "sync"
)

// Task represents a unit of work
type Task interface {
    Execute(ctx context.Context) error
}

// Pool represents a worker pool
type Pool struct {
    workers    int
    taskQueue  chan Task
    resultChan chan Result
    wg         sync.WaitGroup
    ctx        context.Context
    cancel     context.CancelFunc
}

// Result represents task execution result
type Result struct {
    Task  Task
    Error error
}

// NewPool creates a new worker pool
func NewPool(workers int, queueSize int) *Pool {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &Pool{
        workers:    workers,
        taskQueue:  make(chan Task, queueSize),
        resultChan: make(chan Result, queueSize),
        ctx:        ctx,
        cancel:     cancel,
    }
}

// Start starts the worker pool
func (p *Pool) Start() {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker(i)
    }
}

// Stop stops the worker pool gracefully
func (p *Pool) Stop() {
    close(p.taskQueue)
    p.wg.Wait()
    close(p.resultChan)
    p.cancel()
}

// Submit submits a task to the pool
func (p *Pool) Submit(task Task) error {
    select {
    case p.taskQueue <- task:
        return nil
    case <-p.ctx.Done():
        return fmt.Errorf("pool is shutting down")
    }
}

// Results returns the result channel
func (p *Pool) Results() <-chan Result {
    return p.resultChan
}

// worker is the worker goroutine
func (p *Pool) worker(id int) {
    defer p.wg.Done()
    
    for {
        select {
        case task, ok := <-p.taskQueue:
            if !ok {
                return // Channel closed, worker should exit
            }
            
            // Execute task
            err := task.Execute(p.ctx)
            
            // Send result
            select {
            case p.resultChan <- Result{Task: task, Error: err}:
            case <-p.ctx.Done():
                return
            }
            
        case <-p.ctx.Done():
            return
        }
    }
}

// Example usage:
// pool := NewPool(10, 100)
// pool.Start()
// defer pool.Stop()
//
// // Submit tasks
// for i := 0; i < 50; i++ {
//     task := &MyTask{ID: i}
//     pool.Submit(task)
// }
//
// // Process results
// go func() {
//     for result := range pool.Results() {
//         if result.Error != nil {
//             log.Printf("Task failed: %v", result.Error)
//         }
//     }
// }()
```

### HTTP Handler Pattern

```go
// internal/interfaces/http/user_handler.go - HTTP handlers
package http

import (
    "encoding/json"
    "net/http"
    
    "github.com/gorilla/mux"
    "github.com/google/uuid"
    
    "myapp/internal/domain/user"
    "myapp/pkg/errors"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
    userService *user.Service
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *user.Service) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req struct {
        FirstName string `json:"firstName" validate:"required"`
        LastName  string `json:"lastName" validate:"required"`
        Email     string `json:"email" validate:"required,email"`
    }
    
    // Parse request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.writeError(w, errors.NewValidationError("invalid request body", nil))
        return
    }
    
    // Create user
    u, err := h.userService.CreateUser(r.Context(), req.FirstName, req.LastName, req.Email)
    if err != nil {
        h.writeError(w, err)
        return
    }
    
    // Return response
    response := struct {
        ID        string `json:"id"`
        FirstName string `json:"firstName"`
        LastName  string `json:"lastName"`
        Email     string `json:"email"`
    }{
        ID:        u.ID().String(),
        FirstName: u.FirstName(),
        LastName:  u.LastName(),
        Email:     u.Email(),
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

// GetUser handles GET /users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    
    id, err := uuid.Parse(idStr)
    if err != nil {
        h.writeError(w, errors.NewValidationError("invalid user ID", nil))
        return
    }
    
    u, err := h.userService.GetUser(r.Context(), id)
    if err != nil {
        h.writeError(w, err)
        return
    }
    
    response := struct {
        ID        string `json:"id"`
        FirstName string `json:"firstName"`
        LastName  string `json:"lastName"`
        Email     string `json:"email"`
    }{
        ID:        u.ID().String(),
        FirstName: u.FirstName(),
        LastName:  u.LastName(),
        Email:     u.Email(),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// writeError writes error response
func (h *UserHandler) writeError(w http.ResponseWriter, err error) {
    var appErr *errors.AppError
    if errors.As(err, &appErr) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(appErr.HTTPStatus())
        json.NewEncoder(w).Encode(appErr)
        return
    }
    
    // Unknown error
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(w).Encode(map[string]string{
        "error": "internal server error",
    })
}
```

## 🚨 Rule Enforcement

### Validation Checklist

```bash
# Run quality checks
go build ./...
go test ./...
go vet ./...
golangci-lint run

# Check for race conditions
go test -race ./...

# Check module dependencies
go mod tidy
go mod verify

# Security scan
gosec ./...
```

### CI/CD Integration

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.23, 1.24]
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}
    
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
    
    - name: Download dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: Run vet
      run: go vet ./...
    
    - name: Run golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
    
    - name: Run gosec
      uses: securecodewarrior/github-action-gosec@master
      with:
        args: './...'
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
```

### Code Quality Configuration

```yaml
# .golangci.yml
run:
  timeout: 5m
  issues-exit-code: 1
  tests: true

output:
  format: colored-line-number
  print-issued-lines: true
  print-linter-name: true

linters-settings:
  govet:
    check-shadowing: true
  golint:
    min-confidence: 0.8
  gocyclo:
    min-complexity: 15
  maligned:
    suggest-new: true
  dupl:
    threshold: 100
  goconst:
    min-len: 3
    min-occurrences: 3

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - typecheck
    - unused
    - varcheck
    - structcheck
    - deadcode
    - gocyclo
    - gofmt
    - goimports
    - golint
    - gosec
    - interfacer
    - unconvert
    - dupl
    - goconst
    - gocognit
    - bodyclose
    - misspell
    - lll
    - unparam
    - nakedret
    - prealloc
    - scopelint
    - gocritic

issues:
  exclude-use-default: false
  exclude:
    - "should have comment or be unexported"
    - "comment on exported .* should be of the form"
```

## 📖 Rule Reference

| Rule File | Purpose | When Applied |
|-----------|---------|--------------|
| `core/main.mdc` | Main rule coordinator | Always |
| `core/complexity-detection.mdc` | Project analysis | Project initialization |
| `core/module-management.mdc` | Module organization | All projects |
| `modules/error-handling/structured-errors.mdc` | Custom error types | Library packages |
| `modules/error-handling/error-wrapping.mdc` | Error context | All error handling |
| `modules/serialization/json-patterns.mdc` | JSON serialization | API responses |
| `modules/builder-patterns/functional-options.mdc` | Complex constructors | ≥3 parameters |
| `modules/database/repository-patterns.mdc` | Database patterns | Database usage |
| `modules/testing/test-organization.mdc` | Testing patterns | All code |
| `modules/concurrency/modern-patterns.mdc` | High-performance concurrency | Concurrent operations |
| `workflows/development-workflow.mdc` | Complete development process | Always |

## 🎯 Advanced Patterns

### Configuration Management

```go
// internal/config/config.go - Environment-based configuration
package config

import (
    "os"
    "strconv"
    "time"
)

// Config holds application configuration
type Config struct {
    Server   ServerConfig   `json:"server"`
    Database DatabaseConfig `json:"database"`
    Redis    RedisConfig    `json:"redis"`
    Logger   LoggerConfig   `json:"logger"`
}

type ServerConfig struct {
    Host         string        `json:"host"`
    Port         int           `json:"port"`
    ReadTimeout  time.Duration `json:"readTimeout"`
    WriteTimeout time.Duration `json:"writeTimeout"`
    IdleTimeout  time.Duration `json:"idleTimeout"`
}

type DatabaseConfig struct {
    Host            string        `json:"host"`
    Port            int           `json:"port"`
    Username        string        `json:"username"`
    Password        string        `json:"-"` // Excluded from JSON
    Database        string        `json:"database"`
    MaxOpenConns    int           `json:"maxOpenConns"`
    MaxIdleConns    int           `json:"maxIdleConns"`
    ConnMaxLifetime time.Duration `json:"connMaxLifetime"`
}

type RedisConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Password string `json:"-"`
    DB       int    `json:"db"`
}

type LoggerConfig struct {
    Level  string `json:"level"`
    Format string `json:"format"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
    return &Config{
        Server: ServerConfig{
            Host:         getEnv("SERVER_HOST", "localhost"),
            Port:         getEnvInt("SERVER_PORT", 8080),
            ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
            WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
            IdleTimeout:  getEnvDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
        },
        Database: DatabaseConfig{
            Host:            getEnv("DB_HOST", "localhost"),
            Port:            getEnvInt("DB_PORT", 5432),
            Username:        getEnv("DB_USERNAME", "postgres"),
            Password:        getEnv("DB_PASSWORD", ""),
            Database:        getEnv("DB_DATABASE", "myapp"),
            MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
            MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 25),
            ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
        },
        Redis: RedisConfig{
            Host:     getEnv("REDIS_HOST", "localhost"),
            Port:     getEnvInt("REDIS_PORT", 6379),
            Password: getEnv("REDIS_PASSWORD", ""),
            DB:       getEnvInt("REDIS_DB", 0),
        },
        Logger: LoggerConfig{
            Level:  getEnv("LOG_LEVEL", "info"),
            Format: getEnv("LOG_FORMAT", "json"),
        },
    }
}

// Helper functions
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}
```

### Structured Logging

```go
// pkg/logger/logger.go - Structured logging with context
package logger

import (
    "context"
    "os"
    
    "github.com/sirupsen/logrus"
)

type contextKey string

const (
    loggerKey contextKey = "logger"
    fieldsKey contextKey = "fields"
)

// Logger interface for structured logging
type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Fatal(msg string, fields ...Field)
    
    WithFields(fields ...Field) Logger
    WithContext(ctx context.Context) Logger
}

// Field represents a log field
type Field struct {
    Key   string
    Value interface{}
}

// F creates a new field (shorthand)
func F(key string, value interface{}) Field {
    return Field{Key: key, Value: value}
}

// logrusLogger implements Logger interface using logrus
type logrusLogger struct {
    entry *logrus.Entry
}

// NewLogger creates a new logger instance
func NewLogger(level string, format string) Logger {
    log := logrus.New()
    
    // Set log level
    switch level {
    case "debug":
        log.SetLevel(logrus.DebugLevel)
    case "info":
        log.SetLevel(logrus.InfoLevel)
    case "warn":
        log.SetLevel(logrus.WarnLevel)
    case "error":
        log.SetLevel(logrus.ErrorLevel)
    default:
        log.SetLevel(logrus.InfoLevel)
    }
    
    // Set formatter
    if format == "json" {
        log.SetFormatter(&logrus.JSONFormatter{})
    } else {
        log.SetFormatter(&logrus.TextFormatter{
            FullTimestamp: true,
        })
    }
    
    log.SetOutput(os.Stdout)
    
    return &logrusLogger{
        entry: logrus.NewEntry(log),
    }
}

// Implementation of Logger interface
func (l *logrusLogger) Debug(msg string, fields ...Field) {
    l.entry.WithFields(l.convertFields(fields...)).Debug(msg)
}

func (l *logrusLogger) Info(msg string, fields ...Field) {
    l.entry.WithFields(l.convertFields(fields...)).Info(msg)
}

func (l *logrusLogger) Warn(msg string, fields ...Field) {
    l.entry.WithFields(l.convertFields(fields...)).Warn(msg)
}

func (l *logrusLogger) Error(msg string, fields ...Field) {
    l.entry.WithFields(l.convertFields(fields...)).Error(msg)
}

func (l *logrusLogger) Fatal(msg string, fields ...Field) {
    l.entry.WithFields(l.convertFields(fields...)).Fatal(msg)
}

func (l *logrusLogger) WithFields(fields ...Field) Logger {
    return &logrusLogger{
        entry: l.entry.WithFields(l.convertFields(fields...)),
    }
}

func (l *logrusLogger) WithContext(ctx context.Context) Logger {
    // Add context fields if present
    if fields := GetFieldsFromContext(ctx); len(fields) > 0 {
        return &logrusLogger{
            entry: l.entry.WithFields(l.convertFields(fields...)),
        }
    }
    return l
}

func (l *logrusLogger) convertFields(fields ...Field) logrus.Fields {
    logrusFields := make(logrus.Fields)
    for _, field := range fields {
        logrusFields[field.Key] = field.Value
    }
    return logrusFields
}

// Context helpers
func WithLogger(ctx context.Context, logger Logger) context.Context {
    return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) Logger {
    if logger, ok := ctx.Value(loggerKey).(Logger); ok {
        return logger
    }
    return NewLogger("info", "json") // Default logger
}

func WithFields(ctx context.Context, fields ...Field) context.Context {
    existing := GetFieldsFromContext(ctx)
    allFields := append(existing, fields...)
    return context.WithValue(ctx, fieldsKey, allFields)
}

func GetFieldsFromContext(ctx context.Context) []Field {
    if fields, ok := ctx.Value(fieldsKey).([]Field); ok {
        return fields
    }
    return nil
}
```

### Middleware Patterns

```go
// internal/interfaces/http/middleware/middleware.go - HTTP middleware
package middleware

import (
    "context"
    "net/http"
    "time"
    
    "github.com/google/uuid"
    
    "myapp/pkg/logger"
)

// RequestID middleware adds request ID to context
func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := uuid.New().String()
        
        // Add to context
        ctx := context.WithValue(r.Context(), "request_id", requestID)
        
        // Add to response header
        w.Header().Set("X-Request-ID", requestID)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Logging middleware logs HTTP requests
func Logging(log logger.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Add request fields to context
            ctx := logger.WithFields(r.Context(),
                logger.F("method", r.Method),
                logger.F("path", r.URL.Path),
                logger.F("remote_addr", r.RemoteAddr),
                logger.F("user_agent", r.UserAgent()),
            )
            
            // Create response writer wrapper to capture status code
            wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
            
            // Process request
            next.ServeHTTP(wrapped, r.WithContext(ctx))
            
            // Log request completion
            duration := time.Since(start)
            log.WithContext(ctx).Info("HTTP request completed",
                logger.F("status_code", wrapped.statusCode),
                logger.F("duration_ms", duration.Milliseconds()),
            )
        })
    }
}

// responseWriter wrapper to capture status code
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

// Recovery middleware recovers from panics
func Recovery(log logger.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if err := recover(); err != nil {
                    log.WithContext(r.Context()).Error("Panic recovered",
                        logger.F("error", err),
                        logger.F("path", r.URL.Path),
                    )
                    
                    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                }
            }()
            
            next.ServeHTTP(w, r)
        })
    }
}

// CORS middleware handles cross-origin requests
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            
            // Check if origin is allowed
            allowed := false
            for _, allowedOrigin := range allowedOrigins {
                if allowedOrigin == "*" || allowedOrigin == origin {
                    allowed = true
                    break
                }
            }
            
            if allowed {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
                w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
                w.Header().Set("Access-Control-Allow-Credentials", "true")
            }
            
            // Handle preflight requests
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

### Main Application Setup

```go
// cmd/api-server/main.go - Application entry point
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    
    "myapp/internal/config"
    "myapp/internal/domain/user"
    "myapp/internal/infrastructure/database"
    httpHandler "myapp/internal/interfaces/http"
    "myapp/internal/interfaces/http/middleware"
    "myapp/pkg/logger"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        panic(fmt.Sprintf("Failed to load config: %v", err))
    }
    
    // Initialize logger
    log := logger.NewLogger(cfg.Logger.Level, cfg.Logger.Format)
    
    // Initialize database
    db, err := initDatabase(cfg.Database)
    if err != nil {
        log.Fatal("Failed to initialize database", logger.F("error", err))
    }
    defer db.Close()
    
    // Initialize repositories
    userRepo := database.NewUserRepository(db)
    
    // Initialize services
    userService := user.NewService(userRepo)
    
    // Initialize handlers
    userHandler := httpHandler.NewUserHandler(userService)
    
    // Setup routes
    router := setupRoutes(userHandler, log)
    
    // Setup server
    server := &http.Server{
        Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
        Handler:      router,
        ReadTimeout:  cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
        IdleTimeout:  cfg.Server.IdleTimeout,
    }
    
    // Start server in goroutine
    go func() {
        log.Info("Starting server", 
            logger.F("host", cfg.Server.Host),
            logger.F("port", cfg.Server.Port),
        )
        
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("Server failed to start", logger.F("error", err))
        }
    }()
    
    // Wait for interrupt signal to gracefully shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Info("Shutting down server...")
    
    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown", logger.F("error", err))
    }
    
    log.Info("Server exited")
}

func initDatabase(cfg config.DatabaseConfig) (*sqlx.DB, error) {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
    
    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        return nil, err
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(cfg.MaxOpenConns)
    db.SetMaxIdleConns(cfg.MaxIdleConns)
    db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
    
    // Test connection
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    return db, nil
}

func setupRoutes(userHandler *httpHandler.UserHandler, log logger.Logger) *mux.Router {
    router := mux.NewRouter()
    
    // Apply global middleware
    router.Use(middleware.RequestID)
    router.Use(middleware.Logging(log))
    router.Use(middleware.Recovery(log))
    router.Use(middleware.CORS([]string{"*"})) // Configure properly for production
    
    // API routes
    api := router.PathPrefix("/api/v1").Subrouter()
    
    // User routes
    api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
    api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
    
    // Health check
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }).Methods("GET")
    
    return router
}
```

## 🎯 Benefits

### 🔄 Consistency
- Uniform error handling patterns across all packages
- Consistent JSON output format for APIs
- Standardized dependency management with go.mod

### 📈 Scalability
- Clean architecture with domain-driven design
- Interface-based dependency injection
- Modular structure for large teams

### 🛡️ Reliability
- Comprehensive error handling and wrapping
- Proper context usage for cancellation
- Structured logging for debugging

### 👥 Developer Experience
- Clear architectural guidelines
- Comprehensive testing patterns
- Automated quality enforcement

## 🚀 Getting Started

1. **Clone/copy the rules** into your project's `.cursor/go/` directory
2. **Initialize your project** with `go mod init your-project-name`
3. **Set up the directory structure** following the complexity assessment
4. **Apply the appropriate patterns** based on your project type
5. **Run quality checks** to ensure compliance
6. **Integrate CI/CD** for continuous validation

## 📝 Additional Resources

### Recommended Dependencies

```bash
# Core dependencies
go get github.com/gorilla/mux              # HTTP router
go get github.com/jmoiron/sqlx             # SQL extensions
go get github.com/lib/pq                   # PostgreSQL driver
go get github.com/google/uuid              # UUID generation
go get github.com/sirupsen/logrus          # Structured logging

# Testing dependencies
go get github.com/stretchr/testify         # Testing toolkit
go get github.com/golang/mock              # Mock generation
```

### Project Template

```bash
# Create new project
mkdir my-go-app
cd my-go-app
go mod init github.com/username/my-go-app

# Create directory structure
mkdir -p cmd/api-server
mkdir -p internal/{app,domain,infrastructure}
mkdir -p pkg/{client,models}
mkdir -p api/openapi
mkdir -p configs
mkdir -p scripts
mkdir -p docs

# Copy rule files
cp -r .cursor/go/ ./
```

The rule system will guide you through building robust, maintainable Go applications that follow best practices and scale effectively with modern Go 1.23+ features.
        