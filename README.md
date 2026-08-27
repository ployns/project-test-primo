# Product API - Clean Architecture Go Service

A production-ready REST API service built with Golang following clean architecture principles, dependency injection, comprehensive testing, and PostgreSQL database integration.

## 🏗️ Architecture

This project follows **Clean Architecture** with clear separation of concerns:

```
cmd/
├── main.go                 # Application entry point

internal/
├── domain/                 # Business entities & rules
│   └── product.go
├── usecase/               # Business logic (orchestrators)
│   ├── product.go         # Use case implementations
│   └── errors.go          # Domain errors
├── repository/            # Data access layer (interfaces & implementations)
│   ├── repository.go      # Repository implementation
│   └── repository_test.go # Integration tests
├── handler/               # HTTP handlers (Gin framework)
│   ├── product.go         # Request/response handling
│   └── product_test.go    # E2E handler tests
├── db/                    # Database connection & schema
│   └── db.go
└── di/                    # Dependency injection setup
    └── wire.go            # Wire configuration

docs/                      # Swagger documentation
migrations/               # Database migrations

tests/
├── e2e_test.go           # End-to-end tests
└── integration_test.go   # Integration test suite
```

## 📋 Testing Strategy

### 1. **Unit Tests (Domain & Use Case)**
Located in `internal/usecase/product_test.go`
- Tests business logic without external dependencies
- Uses mock repositories
- Tests validation and error handling
- Coverage: CreateProduct, UpdateProduct, GetProduct, DeleteProduct

```bash
go test -v ./internal/usecase/...
```

### 2. **Repository Tests (Integration)**
Located in `internal/repository/repository_test.go`
- Tests database operations with real PostgreSQL
- Tests CRUD operations
- Tests transaction handling
- Coverage: Create, Read, Update, Delete

```bash
go test -v ./internal/repository/...
```

### 3. **Handler Tests (E2E within service)**
Located in `internal/handler/product_test.go`
- Tests HTTP handlers without external calls
- Tests request validation and response format
- Tests error scenarios
- Coverage: All endpoints and error cases

```bash
go test -v ./internal/handler/...
```

### 4. **Full Test Suite**
```bash
go test -v ./...
```

## 🚀 Getting Started

### Prerequisites
- Go 1.27.0+
- PostgreSQL 12+
- Docker & Docker Compose (optional)

### Installation

1. **Clone the repository**
```bash
cd c:\Users\usEr\p\test
```

2. **Install dependencies**
```bash
make install
# or
go mod download
```

3. **Set up environment variables**
```bash
cp .env.example .env
```

Edit `.env` with your database credentials:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=product_db
DB_SSLMODE=disable
PORT=8080
```

4. **Start PostgreSQL (using Docker)**
```bash
make docker-up
```

Or use your existing PostgreSQL installation.

5. **Run the application**
```bash
make run
# or
go run cmd/main.go
```

The server will start at `http://localhost:8080`

## 📚 API Documentation

Swagger documentation is available at: **`http://localhost:8080/api-docs/index.html`**

### Endpoints

#### 1. Create Product
```http
POST /product
Content-Type: application/json

{
  "name": "Product Name",
  "description": "Product description (optional)",
  "price": 99.99,
  "sale_price": 79.99  (optional)
}
```

**Response:**
```json
{
  "successful": true,
  "error_code": "",
  "data": {
    "data1": 1,
    "data2": "Product Name",
    "price": 99.99,
    "sale_price": 79.99
  }
}
```

#### 2. Get Product
```http
GET /product/{id}
```

**Response:**
```json
{
  "successful": true,
  "error_code": "",
  "data": {
    "data1": 1,
    "data2": "Product Name",
    "price": 99.99,
    "sale_price": 79.99
  }
}
```

#### 3. Update Product (Partial Update)
```http
PATCH /product/{id}
Content-Type: application/json

{
  "name": "Updated Name",
  "price": 149.99
}
```

Only included fields will be updated. Response:
```json
{
  "successful": true,
  "error_code": ""
}
```

#### 4. Delete Product
```http
DELETE /product/{id}
```

**Response:**
```json
{
  "successful": true,
  "error_code": ""
}
```

## 🔧 Dependency Injection

This project uses **Google Wire** for compile-time dependency injection.

### Wire Configuration
Located in `internal/di/wire.go`:

```go
func InitializeProductHandler(postgres *sql.DB) *handler.ProductHandler
func InitializeDatabase() (*sql.DB, error)
func InitializeApp() (*sql.DB, *handler.ProductHandler, error)
```

### Generate Wire Code
```bash
make wire
# or
wire ./internal/di
```

**Benefits:**
- Type-safe dependency injection
- Compile-time verification
- No reflection overhead
- Clear dependency graph

## 🗄️ Database Schema

PostgreSQL schema is automatically created on startup:

```sql
CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  price DOUBLE PRECISION NOT NULL,
  sale_price DOUBLE PRECISION,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_id ON products(id);
```

## 📝 Makefile Commands

```bash
make install         # Install dependencies
make build          # Build the application
make run            # Run the application
make test           # Run all tests
make test-unit      # Run unit tests only
make test-repo      # Run repository integration tests
make swagger        # Generate Swagger documentation
make wire           # Generate Wire DI code
make clean          # Clean build artifacts
make docker-up      # Start PostgreSQL container
make docker-down    # Stop PostgreSQL container
make lint           # Run linter
make fmt            # Format code
make vet            # Run go vet
```

## 🧪 Running Tests

### Run All Tests
```bash
go test -v ./...
```

### Run Specific Test Suite
```bash
# Unit tests (mocked)
go test -v ./internal/usecase/...
go test -v ./internal/handler/...

# Integration tests (with database)
go test -v ./internal/repository/...
```

### Run Tests with Coverage
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📦 Dependencies

Core dependencies in `go.mod`:

- **gin-gonic/gin** - HTTP framework
- **lib/pq** - PostgreSQL driver
- **swaggo/swag** - Swagger code generation
- **google/wire** - Dependency injection
- **stretchr/testify** - Testing assertions
- **joho/godotenv** - Environment variables

## 🏭 Clean Architecture Layers

### 1. **Domain Layer** (`internal/domain/`)
- Pure business entities and rules
- No dependencies on other layers
- Contains: `Product`, request/response models

### 2. **Use Case Layer** (`internal/usecase/`)
- Business logic orchestration
- Implements repository interfaces (Dependency Inversion)
- Contains: `ProductUseCase` with CRUD operations
- Independent of HTTP framework and database

### 3. **Repository Layer** (`internal/repository/`)
- Data access abstraction
- Implements repository interfaces from use case layer
- Contains: Database operations and SQL queries
- Can be easily swapped with other implementations (e.g., MongoDB)

### 4. **Handler Layer** (`internal/handler/`)
- HTTP request/response handling
- Framework-specific code (Gin)
- Delegates to use cases
- Responsible for validation and response formatting

### 5. **Infrastructure Layer** (`internal/db/`, `internal/di/`)
- Database connections
- Dependency injection wiring
- Configuration management

## 🔐 Dependency Rule

```
Entities → UseCases → Interfaces (Repository) ← Implementations
  ↑         ↑
All inner layers point inward
Outer layers depend on inner layers
```

No outer layer should know about inner layers directly.

## 🚦 Error Handling

Error codes used in API responses:
- `INVALID_REQUEST` - Malformed request body
- `CREATE_FAILED` - Product creation failed
- `UPDATE_FAILED` - Product update failed
- `DELETE_FAILED` - Product deletion failed
- `INVALID_ID` - Invalid product ID format
- `NOT_FOUND` - Product not found
- `INTERNAL_ERROR` - Server error

## 🔄 Workflow

```
HTTP Request
    ↓
Handler (validates request format)
    ↓
UseCase (validates business logic)
    ↓
Repository (database operations)
    ↓
Database
    ↓
(Response flows back through the same layers)
```

## 📈 Scalability Considerations

1. **Database**: Can switch to connection pooling
2. **Repository**: Easily swap PostgreSQL for MongoDB, MySQL, etc.
3. **Caching**: Add caching layer between UseCase and Repository
4. **Messaging**: Add event publishing in UseCase layer
5. **Logging**: Add structured logging middleware
6. **Metrics**: Add Prometheus metrics collection

## 🛠️ Development

### Code Format
```bash
make fmt  # Format code
make vet  # Run go vet
```

### Regenerate Swagger
```bash
make swagger
```

### Regenerate Wire
```bash
make wire
```

## 📋 Checklist for Production

- [ ] Add structured logging (zap, logrus)
- [ ] Add request/response logging middleware
- [ ] Add error tracking (Sentry)
- [ ] Add metrics (Prometheus)
- [ ] Add health check endpoint
- [ ] Add graceful shutdown
- [ ] Add request timeout configuration
- [ ] Add CORS configuration
- [ ] Add rate limiting
- [ ] Add authentication/authorization
- [ ] Add database migrations (Flyway, Migrate)
- [ ] Add CI/CD pipeline
- [ ] Add Docker containerization
- [ ] Add Kubernetes manifests
- [ ] Add API versioning

## 📞 Support

For issues or questions, refer to the test files for usage examples.

---

**Author:** Product Team  
**Last Updated:** 2026-08-25  
**License:** MIT
