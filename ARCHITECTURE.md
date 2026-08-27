# Clean Architecture & Dependency Injection Guide

## Overview

This project implements **Clean Architecture** principles to create a maintainable, testable, and scalable Go service. The architecture is organized in concentric circles, where dependencies point inward.

## Architecture Layers

### 1. Domain Layer (innermost)
**Location:** `internal/domain/`

Contains pure business entities and rules, independent of any framework or external dependency.

```go
// Entities
type Product struct {
    ID          int64
    Name        string
    Description *string
    Price       float64
    SalePrice   *float64
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Request/Response DTOs
type CreateProductRequest struct {
    Name      string  `json:"name"`
    Price     float64 `json:"price"`
    ...
}
```

**Characteristics:**
- No external dependencies
- Pure Go - no imports except standard library
- Business logic and validation rules
- Request/Response models (DTOs)

---

### 2. Use Case Layer
**Location:** `internal/usecase/`

Implements business logic orchestration. This layer:
- Coordinates domain entities
- Implements business workflows
- Uses repository **interfaces** (not implementations)
- Defines its own error types

```go
type ProductUseCase interface {
    CreateProduct(ctx context.Context, req *CreateProductRequest) (*Product, error)
    GetProduct(ctx context.Context, id int64) (*Product, error)
    UpdateProduct(ctx context.Context, id int64, req *UpdateProductRequest) error
    DeleteProduct(ctx context.Context, id int64) error
}

type productUseCase struct {
    repo repository.ProductRepository  // Depends on interface, not implementation
}

func (uc *productUseCase) CreateProduct(ctx context.Context, req *CreateProductRequest) (*Product, error) {
    // Validate business rules
    if req.Name == "" {
        return nil, ErrInvalidProductName
    }
    if req.Price <= 0 {
        return nil, ErrInvalidPrice
    }
    
    // Call repository
    return uc.repo.Create(ctx, product)
}
```

**Characteristics:**
- Depends on interfaces (repository), not implementations
- Implements business rules validation
- Coordinates entity operations
- No knowledge of HTTP or database

---

### 3. Interface Adapters Layer
**Location:** `internal/handler/`, `internal/repository/`

Converts data between use cases and external systems.

#### Handlers (Input Adapters)
```go
type ProductHandler struct {
    uc usecase.ProductUseCase
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
    // Parse HTTP request
    var req domain.CreateProductRequest
    c.ShouldBindJSON(&req)
    
    // Call use case
    product, err := h.uc.CreateProduct(c.Request.Context(), &req)
    
    // Format HTTP response
    c.JSON(http.StatusOK, domain.Response{...})
}
```

#### Repositories (Output Adapters)
```go
type ProductRepository interface {
    Create(ctx context.Context, product *Product) (*Product, error)
    GetByID(ctx context.Context, id int64) (*Product, error)
    Update(ctx context.Context, id int64, product *Product) error
    Delete(ctx context.Context, id int64) error
}

type productRepository struct {
    db *sql.DB
}

func (r *productRepository) Create(ctx context.Context, product *Product) (*Product, error) {
    // Execute SQL
    row := r.db.QueryRowContext(ctx, query, ...)
    // Map to entity
    err := row.Scan(&product.ID, ...)
    return product, nil
}
```

---

### 4. Framework & Drivers Layer (outermost)
**Location:** `cmd/`, `internal/db/`, `internal/di/`

External libraries and frameworks.

```go
// cmd/main.go - Entry point
func main() {
    database, err := db.NewPostgresDB()
    
    // Setup dependency injection
    productRepo := repository.NewProductRepository(database)
    productUC := usecase.NewProductUseCase(productRepo)
    productHandler := handler.NewProductHandler(productUC)
    
    // Setup HTTP framework
    router := gin.Default()
    productHandler.RegisterRoutes(router)
    
    router.Run(":8080")
}
```

---

## Dependency Inversion Principle

The key to clean architecture: **Depend on abstractions, not concretions**.

### ❌ Wrong Approach (Tight Coupling)
```go
// UseCase depends on concrete implementation
type productUseCase struct {
    repo *productRepository  // Direct dependency on concrete type
}
```

### ✅ Correct Approach (Loose Coupling)
```go
// UseCase depends on interface
type productUseCase struct {
    repo repository.ProductRepository  // Interface dependency
}

// Interface is defined in the use case layer
type ProductRepository interface {
    Create(ctx context.Context, product *Product) (*Product, error)
    // ...
}
```

**Benefits:**
- Easy to test (mock the interface)
- Easy to swap implementations
- Follows SOLID principles
- Lower coupling between layers

---

## Dependency Injection with Wire

This project uses **Google Wire** for compile-time dependency injection.

### Wire Configuration
**Location:** `internal/di/wire.go`

```go
// +build wireinject

package di

func InitializeProductHandler(postgres *sql.DB) *handler.ProductHandler {
    wire.Build(
        repository.NewProductRepository,
        usecase.NewProductUseCase,
        handler.NewProductHandler,
    )
    return nil
}
```

### How Wire Works

1. **Declares** constructors that create dependencies
2. **Analyzes** function signatures at compile time
3. **Generates** `wire_gen.go` with complete dependency graph
4. **Verifies** all dependencies are satisfied

### Generate Wire Code
```bash
wire ./internal/di
```

**Advantages over runtime DI:**
- Compile-time verification (catch errors early)
- No reflection overhead
- Type-safe
- Clear dependency graph

---

## Data Flow

### Incoming Request
```
HTTP Request
    ↓
Handler (gin.Context)
    ├─ Parse: JSON → CreateProductRequest DTO
    ├─ Validate: Format and constraints
    └─ Call UseCase
        ↓
UseCase
    ├─ Validate: Business rules
    ├─ Call: Repository interface
    └─ Return: Product entity
        ↓
Repository
    ├─ Execute: SQL query
    ├─ Map: DB rows → Product entity
    └─ Return: Product
        ↓
UseCase (receives Product)
    └─ Return: Product
        ↓
Handler (receives Product)
    ├─ Transform: Product → ProductData DTO
    ├─ Format: JSON response
    └─ Write HTTP Response
        ↓
HTTP Response (JSON)
```

### Key Points:
- **Handler** converts HTTP ↔ DTOs
- **UseCase** implements business logic
- **Repository** converts DTOs ↔ Database rows
- Each layer has clear responsibility
- Dependencies flow inward only

---

## Testing Strategy

### 1. Unit Tests (UseCase)
```go
// Tests business logic without database
mockRepo := new(mockProductRepository)
mockRepo.On("Create", ...).Return(product, nil)

uc := NewProductUseCase(mockRepo)
result, err := uc.CreateProduct(ctx, req)
```

**Characteristics:**
- Mocks repository
- Fast execution
- No external dependencies
- Tests business logic in isolation

---

### 2. Integration Tests (Repository)
```go
// Tests database operations
db := sql.Open("postgres", dsn)
repo := NewProductRepository(db)

product, err := repo.Create(ctx, testProduct)
assert.NoError(t, err)
assert.Greater(t, product.ID, int64(0))
```

**Characteristics:**
- Uses real PostgreSQL
- Tests data persistence
- Slower than unit tests
- Tests database layer in isolation

---

### 3. Handler Tests (HTTP)
```go
// Tests HTTP layer without real repository
mockUC := new(mockUseCase)
handler := NewProductHandler(mockUC)

router := gin.New()
handler.RegisterRoutes(router)

req := httptest.NewRequest("POST", "/product", body)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)
```

**Characteristics:**
- Mocks use case
- Tests request parsing and response formatting
- Tests HTTP status codes and errors
- Doesn't touch database

---

### 4. E2E Tests
```go
// Tests complete flow from HTTP to database
db := setupTestDB()
router := setupRouter(db)

// Create request, execute, verify response
createReq := httptest.NewRequest("POST", "/product", body)
w := httptest.NewRecorder()
router.ServeHTTP(w, createReq)

// Verify data in database
product, _ := repo.GetByID(ctx, 1)
assert.Equal(t, "Product Name", product.Name)
```

**Characteristics:**
- Tests complete request-response cycle
- Uses real database
- Most realistic testing
- Slowest but most comprehensive

---

## Error Handling

### Domain Errors
Defined in `internal/usecase/errors.go`

```go
var (
    ErrInvalidProductID   = errors.New("invalid product id")
    ErrInvalidProductName = errors.New("product name cannot be empty")
    ErrInvalidPrice       = errors.New("price must be greater than 0")
)
```

### Error Flow
```
UseCase validation error
    ↓
Returns: (nil, ErrInvalidProductName)
    ↓
Handler catches error
    ↓
Maps to: ErrorCode: "INVALID_REQUEST"
    ↓
Returns: HTTP 400
```

---

## Scalability Considerations

### Caching Layer
```go
type cachedProductRepository struct {
    repo  repository.ProductRepository
    cache Cache
}

func (r *cachedProductRepository) GetByID(ctx context.Context, id int64) (*Product, error) {
    if product, ok := r.cache.Get(id); ok {
        return product, nil
    }
    product, err := r.repo.GetByID(ctx, id)
    r.cache.Set(id, product)
    return product, err
}
```

### Event Publishing
```go
type eventPublishingUseCase struct {
    repo      repository.ProductRepository
    publisher events.Publisher
}

func (uc *eventPublishingUseCase) CreateProduct(...) (*Product, error) {
    product, err := uc.repo.Create(...)
    uc.publisher.Publish("product.created", product)
    return product, err
}
```

### Logging
```go
type loggingUseCase struct {
    uc     ProductUseCase
    logger Logger
}

func (l *loggingUseCase) CreateProduct(ctx context.Context, req *CreateProductRequest) (*Product, error) {
    l.logger.Info("Creating product", "name", req.Name)
    product, err := l.uc.CreateProduct(ctx, req)
    l.logger.Info("Product created", "id", product.ID)
    return product, err
}
```

---

## SOLID Principles

### S - Single Responsibility
- Each layer has one reason to change
- Handlers: HTTP representation
- UseCase: Business logic
- Repository: Data persistence

### O - Open/Closed
- Open for extension: Add new features via new implementations
- Closed for modification: Interfaces define contracts

### L - Liskov Substitution
- Repository implementations are interchangeable
- Can swap PostgreSQL for MongoDB without changing UseCase

### I - Interface Segregation
- Interfaces are focused and minimal
- Repository only exposes needed methods

### D - Dependency Inversion
- Depend on abstractions (interfaces)
- Wire injects concrete implementations
- UseCase doesn't know about database driver

---

## Key Takeaways

1. **Layers are concentric**: Dependencies point inward
2. **Interfaces are boundaries**: Define contracts between layers
3. **Wire provides DI**: Compile-time wiring of dependencies
4. **Test each layer independently**: Unit, integration, E2E tests
5. **Business logic stays pure**: No framework knowledge in domain/usecase
6. **External concerns are isolated**: Database, HTTP in adapter layers

This architecture makes the codebase:
- **Testable**: Mock any layer independently
- **Maintainable**: Clear separation of concerns
- **Scalable**: Easy to add new features and layers
- **Flexible**: Easy to swap implementations
