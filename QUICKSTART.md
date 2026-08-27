# Quick Start Guide

Get the Product API running in 5 minutes!

## Option 1: Using Docker Compose (Recommended)

### 1. Start PostgreSQL
```bash
docker-compose up -d
```

### 2. Copy environment file
```bash
cp .env.example .env
```

### 3. Download dependencies
```bash
go mod download
```

### 4. Run the application
```bash
go run cmd/main.go
```

The API will be available at `http://localhost:8080`

### 5. View Swagger API docs
Open in browser: `http://localhost:8080/api-docs/index.html`

---

## Option 2: Using existing PostgreSQL

### 1. Copy environment file and update it
```bash
cp .env.example .env
```

Edit `.env` with your PostgreSQL connection details:
```env
DB_HOST=your_host
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=product_db
DB_SSLMODE=disable
```

### 2. Download dependencies
```bash
go mod download
```

### 3. Run the application
```bash
go run cmd/main.go
```

---

## Testing the API

### Create a Product
```bash
curl -X POST http://localhost:8080/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 1299.99,
    "sale_price": 999.99
  }'
```

### Get a Product
```bash
curl http://localhost:8080/product/1
```

### Update a Product
```bash
curl -X PATCH http://localhost:8080/product/1 \
  -H "Content-Type: application/json" \
  -d '{
    "price": 1399.99
  }'
```

### Delete a Product
```bash
curl -X DELETE http://localhost:8080/product/1
```

---

## Running Tests

### Run all tests
```bash
go test -v ./...
```

### Run specific test suites
```bash
# Unit tests only
go test -v ./internal/usecase/... ./internal/handler/...

# Integration tests (requires database)
go test -v ./internal/repository/...

# E2E tests
go test -v ./tests/...
```

### With coverage
```bash
go test -cover ./...
```

---

## Project Structure

```
cmd/
  └── main.go          # Application entry point

internal/
  ├── domain/          # Business entities
  ├── usecase/         # Business logic
  ├── repository/      # Database access
  ├── handler/         # HTTP handlers
  ├── db/              # Database setup
  └── di/              # Dependency injection

tests/
  └── e2e_test.go      # End-to-end tests

docs/
  └── swagger.go       # Swagger docs
```

---

## Useful Commands

```bash
# Build the app
make build

# Run the app
make run

# Run all tests
make test

# Generate Swagger docs
make swagger

# Start/stop PostgreSQL
make docker-up
make docker-down

# Code formatting
make fmt

# Code linting
make lint
```

---

## Common Issues

### Database connection refused
- Ensure PostgreSQL is running: `docker-compose up -d`
- Check `.env` file has correct credentials
- Verify PostgreSQL is listening on port 5432

### Port 8080 already in use
- Change `PORT` in `.env` file
- Or stop the process using port 8080

### Tests fail with database error
- For integration tests, start PostgreSQL first: `docker-compose up -d`
- Unit tests don't require a database

---

## Next Steps

1. Read [README.md](README.md) for detailed documentation
2. Check test files for usage examples
3. View Swagger UI at `/api-docs` for interactive API testing
4. Explore the clean architecture structure to understand the design

Happy coding! 🚀
