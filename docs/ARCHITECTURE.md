# Hexagonal Architecture - Movie API

## Overview

This project implements **Hexagonal Architecture** (also known as Ports and Adapters Architecture) to create a REST API for movies in Go. This architecture promotes separation of concerns, testability, and framework independence.

> **Note**: For quick start, installation, and API usage examples, see [../README.md](../README.md).

## Project Structure

```
api-server/
├── config/           # Application configuration
├── database/         # Database configuration and migration
├── handler/          # HTTP adapters (Primary Input Ports)
│   ├── movie_handler.go
│   └── routes.go
├── models/           # Domain models and DTOs
├── repository/       # Database adapters (Secondary Output Ports)
│   └── movie_repository.go
├── service/          # Pure business logic (Domain)
│   ├── movie_service.go
│   └── movie_service_test.go
├── utils/            # General utilities
└── main.go          # Entry point and dependency wiring
```

## Hexagonal Architecture

### Visual Overview

```
┌─────────────────────────────────────────────────────────┐
│                    EXTERNAL WORLD                       │
│                                                         │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐  │
│  │   Client    │    │   Client    │    │   Client    │  │
│  │    HTTP     │    │     CLI     │    │   gRPC      │  │
│  └─────────────┘    └─────────────┘    └─────────────┘  │
│         │                   │                   │       │
│         ▼                   ▼                   ▼       │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐  │
│  │   Handler   │    │   CLI App   │    │  gRPC Serv  │  │
│  │  (Gin)      │    │             │    │             │  │
│  └─────────────┘    └─────────────┘    └─────────────┘  │
│         │                   │                   │       │
│         └───────────────────┼───────────────────┘       │
│                             │                           │
│                    ┌────────▼────────┐                  │
│                    │                 │                  │
│                    │   DOMAIN        │                  │
│                    │   (Service)     │                  │
│                    │                 │                  │
│                    └────────┬────────┘                  │
│                             │                           │
│         ┌───────────────────┼───────────────────┐       │
│         │                   │                   │       │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐  │
│  │   GORM DB   │    │   Redis     │    │   File      │  │
│  │             │    │   Cache     │    │   System    │  │
│  └─────────────┘    └─────────────┘    └─────────────┘  │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### 1. Domain (The Central Hexagon)

**Location**: `service/movie_service.go`

The domain contains pure business logic, without external dependencies:

```go
type MovieService interface {
    GetMovie(id uint) (*models.Movie, error)
    CreateMovie(req *models.MovieCreateRequest) (*models.Movie, error)
    // ... other methods
}
```

**Characteristics**:
- ✅ No dependency on web frameworks (Gin, Echo, etc.)
- ✅ No dependency on databases (GORM, SQLx, etc.)
- ✅ Contains business validations
- ✅ Easy to test in isolation

### 2. Primary Input Ports (HTTP Adapters)

**Location**: `handler/movie_handler.go`

Handlers act as adapters between HTTP and the domain:

```go
type MovieHandler struct {
    service service.MovieService  // Depends on domain interface
}

func (h *MovieHandler) Get(c *gin.Context) {
    // 1. Parse HTTP parameters
    // 2. Call service
    // 3. Format HTTP response
}
```

**Characteristics**:
- ✅ Single responsibility: adapt HTTP ↔ Domain
- ✅ Depends on interfaces, not implementations
- ✅ Easy to change web framework

### 3. Secondary Output Ports (Database Adapters)

**Location**: `repository/movie_repository.go`

Repositories abstract data access:

```go
type MovieRepository interface {
    FindByID(id uint) (*models.Movie, error)
    Create(movie *models.Movie) error
    // ... other methods
}

type gormMovieRepository struct {
    db *gorm.DB
}
```

**Characteristics**:
- ✅ Interface independent of implementation
- ✅ Easy to change database
- ✅ Easy to mock for testing

### 4. Dependency Wiring

**Location**: `main.go`

The entry point configures all dependencies:

```go
func main() {
    // 1. Database
    database.InitDB()
    
    // 2. Repository (data layer)
    movieRepo := repository.NewMovieRepository(database.DB)
    
    // 3. Service (business logic)
    movieService := service.NewMovieService(movieRepo)
    
    // 4. Handler (HTTP adapter)
    movieHandler := handler.NewMovieHandler(movieService)
    
    // 5. Configure routes
    handler.SetupRoutes(app, movieHandler)
}
```

## Benefits of this Architecture

### 1. **High Testability**

```go
// Test service without real database
func TestGetMovie(t *testing.T) {
    mockRepo := NewMockMovieRepository()
    service := NewMovieService(mockRepo)
    
    movie, err := service.GetMovie(1)
    // Assertions...
}
```

### 2. **Framework Independence**

Your business logic works the same if you change from Gin to Echo:

```go
// With Gin
func (h *MovieHandler) Get(c *gin.Context) {
    movie, err := h.service.GetMovie(id)
    c.JSON(200, movie)
}

// With Echo (only handler would change)
func (h *MovieHandler) Get(c echo.Context) {
    movie, err := h.service.GetMovie(id)
    c.JSON(200, movie)
}
```

### 3. **Database Independence**

You can change from GORM to SQLx without touching the domain:

```go
// GORM implementation
type gormRepository struct { db *gorm.DB }

// SQLx implementation
type sqlxRepository struct { db *sqlx.DB }

// Test implementation
type mockRepository struct { movies map[uint]*models.Movie }
```

### 4. **Clear Separation of Responsibilities**

- **Handler**: Only handles HTTP
- **Service**: Only business logic
- **Repository**: Only data access

## Request Flow

```
HTTP Request
    ↓
Handler (Parse HTTP)
    ↓
Service (Business logic)
    ↓
Repository (Data access)
    ↓
Database
```

## SOLID Principles Applied

### 1. **Single Responsibility Principle (SRP)**
- Each layer has a single responsibility
- Handlers only handle HTTP
- Services only contain business logic
- Repositories only access data

### 2. **Open/Closed Principle (OCP)**
- You can extend functionality without modifying existing code
- New handlers without touching services
- New repositories without touching services

### 3. **Liskov Substitution Principle (LSP)**
- Any implementation of `MovieRepository` can substitute another
- Any implementation of `MovieService` can substitute another

### 4. **Interface Segregation Principle (ISP)**
- Small and specific interfaces
- `MovieRepository` only movie methods
- `MovieService` only movie logic

### 5. **Dependency Inversion Principle (DIP)**
- Dependencies on abstractions, not implementations
- `Service` depends on `MovieRepository` (interface)
- `Handler` depends on `MovieService` (interface)

## Development Workflow

### Running Tests
```bash
# Run service tests (business logic)
go test ./service

# Run all tests
go test ./...
```

### Code Quality
```bash
# Format code
gofmt -w .

# Run linter
golint ./...

# Run static analysis
go vet ./...
```

## API Endpoints Reference

> **Note**: For detailed usage examples, see [../README.md](../README.md).

- `GET /health` - Health check
- `GET /movies` - List movies with filters
- `GET /movies/:id` - Get movie by ID
- `POST /movies` - Create new movie
- `PUT /movies/:id` - Update movie
- `DELETE /movies/:id` - Delete movie
- `GET /movies/search?title=...` - Search by title
- `GET /movies/top-rated` - Top rated movies
- `GET /genres/:id/movies` - Movies by genre
- `GET /directors/:id/movies` - Movies by director
- `GET /actors/:id/movies` - Movies by actor

## Next Steps

1. **Add more business validations** in the service
2. **Implement authentication and authorization**
3. **Add structured logging**
4. **Implement cache with Redis**
5. **Add metrics and monitoring**
6. **Implement rate limiting**
7. **Add Swagger documentation**

## Related Documentation

- **[../README.md](../README.md)** - Quick start and API usage
- **[../.cursorrules](../.cursorrules)** - Project coding standards 