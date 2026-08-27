# Windows Setup Guide

This guide covers setting up and running the project on Windows without `make`.

## 🚀 Quick Start (Windows)

### 1. Prerequisites
- Go 1.27.0+
- Git
- PostgreSQL 12+ (or Docker)

### 2. Clone & Setup
```powershell
cd c:\Users\usEr\p\test
go mod download
cp .env.example .env
```

### 3. Start PostgreSQL

**Option A: Using Docker**
```powershell
docker run --name postgres-product `
  -e POSTGRES_PASSWORD=postgres `
  -e POSTGRES_DB=product_db `
  -p 5432:5432 `
  -d postgres:15
```

**Option B: Using docker-compose**
```powershell
docker-compose up -d
```

**Option C: Using existing PostgreSQL**
- Update `.env` with your connection details
- Make sure the database exists

### 4. Run Application
```powershell
go run cmd/main.go
```

Server will start at `http://localhost:8080`

### 5. Test the API
```powershell
# Create a product
$body = @{
    name = "Test Product"
    price = 99.99
} | ConvertTo-Json

$response = Invoke-WebRequest `
  -Uri "http://localhost:8080/product" `
  -Method POST `
  -ContentType "application/json" `
  -Body $body

$response.Content | ConvertFrom-Json

# View Swagger docs
Start-Process "http://localhost:8080/api-docs/index.html"
```

## 🧪 Running Tests (Windows)

### All Tests
```powershell
go test -v ./...
```

### By Layer
```powershell
# Unit tests (Handler layer)
go test -v ./internal/handler

# Unit tests (UseCase layer)
go test -v ./internal/usecase

# Integration tests (Repository layer)
go test -v ./internal/repository

# E2E tests
go test -v ./tests
```

### Unit Tests Only (No Database)
```powershell
go test -v ./internal/handler ./internal/usecase
```

### With Coverage
```powershell
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Specific Test
```powershell
go test -v ./internal/handler -run TestProductHandler_CreateProduct
```

### Run with Short Timeout
```powershell
go test -v -timeout 30s ./internal/handler
```

## 🛠️ Building (Windows)

### Build Executable
```powershell
go build -o bin\app.exe cmd\main.go
```

### Run Built Executable
```powershell
.\bin\app.exe
```

### Clean Build Artifacts
```powershell
Remove-Item -Path bin -Recurse -Force
go clean
```

## 📦 Dependencies

### Download All Dependencies
```powershell
go mod download
```

### Update Dependencies
```powershell
go mod tidy
```

### Check Dependency Versions
```powershell
go mod graph
go list -m all
```

## 🗄️ Database Management (Windows)

### Start PostgreSQL Container
```powershell
# Using docker run
docker run --name postgres-product `
  -e POSTGRES_PASSWORD=postgres `
  -e POSTGRES_DB=product_db `
  -p 5432:5432 `
  -d postgres:15

# Using docker-compose
docker-compose up -d
```

### Stop PostgreSQL Container
```powershell
# Using docker
docker stop postgres-product
docker rm postgres-product

# Using docker-compose
docker-compose down
```

### Access PostgreSQL
```powershell
# Connect to database
docker exec -it postgres-product psql -U postgres -d product_db

# View tables
\dt

# View schema
\d products

# Exit
\q
```

## 🔄 Development Workflow (Windows)

### Setup One-Time
```powershell
cd c:\Users\usEr\p\test
go mod download
docker-compose up -d
```

### Daily Development
```powershell
# 1. Run application
go run cmd/main.go

# 2. In another PowerShell window, run tests
go test -v ./...

# 3. Test API manually
Invoke-WebRequest http://localhost:8080/api-docs/index.html -UseBasicParsing
```

### After Making Changes
```powershell
# Format code
go fmt ./...

# Run linter
golangci-lint run ./...

# Run tests
go test -v ./...

# Build
go build -o bin\app.exe cmd\main.go
```

## 📝 PowerShell Aliases

Add to PowerShell profile (`$PROFILE`) for convenience:

```powershell
# Go test shortcuts
function Test-All {
    go test -v ./...
}

function Test-Unit {
    go test -v ./internal/handler ./internal/usecase
}

function Test-Integration {
    go test -v ./internal/repository
}

function Test-E2E {
    go test -v ./tests
}

# Build shortcuts
function Build-App {
    go build -o bin\app.exe cmd\main.go
}

function Run-App {
    go run cmd\main.go
}

# Database shortcuts
function Start-DB {
    docker-compose up -d
}

function Stop-DB {
    docker-compose down
}

# Format and clean
function Format-Code {
    go fmt ./...
}

function Clean-Build {
    Remove-Item -Path bin -Recurse -Force -ErrorAction SilentlyContinue
    go clean
}
```

Usage:
```powershell
Test-Unit
Build-App
Run-App
Start-DB
Stop-DB
```

## 🐛 Troubleshooting

### Port Already in Use
```powershell
# Find process using port 8080
netstat -ano | findstr :8080

# Kill process (replace PID)
taskkill /PID <PID> /F

# Or use different port
$env:PORT=8081
go run cmd/main.go
```

### Database Connection Failed
```powershell
# Check if PostgreSQL is running
docker ps -a | findstr postgres

# View database logs
docker logs postgres-product

# Restart database
docker-compose restart

# Check connection string in .env
type .env
```

### Tests Failing with Database Error
```powershell
# Ensure database is running
docker-compose up -d

# Check database exists
docker exec -it postgres-product psql -U postgres -l

# Run only unit tests (no database required)
go test -v ./internal/handler ./internal/usecase
```

### Build Fails with Missing Dependencies
```powershell
# Tidy dependencies
go mod tidy

# Download all
go mod download

# Verify
go mod verify
```

## 📊 Environment Setup Examples

### Development (.env)
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=product_db
DB_SSLMODE=disable
PORT=8080
```

### Testing (.env.test)
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=test_db
DB_SSLMODE=disable
PORT=8080
```

### Production (.env.prod)
```env
DB_HOST=prod-db.example.com
DB_PORT=5432
DB_USER=produser
DB_PASSWORD=securepassword
DB_NAME=product_db
DB_SSLMODE=require
PORT=8080
```

## 🎯 Common Tasks

### Run Single Test Function
```powershell
go test -v -run TestProductHandler_CreateProduct ./internal/handler
```

### Run Tests Matching Pattern
```powershell
go test -v -run "CreateProduct" ./...
```

### Run Tests Excluding Pattern
```powershell
go test -v ./... -run "!E2E"
```

### Parallel Test Execution
```powershell
go test -v -parallel 4 ./...
```

### Verbose Output with Coverage
```powershell
go test -v -cover ./...
```

### Generate Coverage Report
```powershell
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
Start-Process coverage.html
```

## 🔗 Useful Links

- [Go Download](https://golang.org/dl)
- [PostgreSQL Download](https://www.postgresql.org/download/windows)
- [Docker Desktop](https://www.docker.com/products/docker-desktop)
- [Go Testing Package](https://pkg.go.dev/testing)

## 📚 Next Steps

1. Run tests: `go test -v ./...`
2. Start app: `go run cmd/main.go`
3. View API docs: `http://localhost:8080/api-docs`
4. Make changes and test
5. Build for deployment: `go build -o bin\app.exe cmd\main.go`

---

**Windows Go Development Guide Complete!** ✅
