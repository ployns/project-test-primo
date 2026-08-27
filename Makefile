.PHONY: help install build run test clean swagger wire mocks docker-up docker-down

help:
	@echo "Available commands:"
	@echo "  make install        - Install dependencies"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run all tests"
	@echo "  make test-unit      - Run unit tests"
	@echo "  make test-repo      - Run repository (integration) tests"
	@echo "  make mocks          - Generate mocks using mockery"
	@echo "  make swagger        - Generate Swagger documentation"
	@echo "  make wire           - Generate Wire dependency injection"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make docker-up      - Start PostgreSQL container"
	@echo "  make docker-down    - Stop PostgreSQL container"

install:
	go mod download
	go mod tidy

build:
	go build -o bin/app cmd/main.go

run: build
	./bin/app

test:
	go test -v ./...

test-unit:
	go test -v -short ./internal/usecase/... ./internal/handler/...

test-repo:
	go test -v ./internal/repository/...

mocks:
	mockery --all --output=internal/usecase/mocks --outpkg=mocks --case=underscore
	mockery --output=internal/repository/mocks --outpkg=mocks --case=underscore \
		-r internal/repository

swagger:
	swag init -g cmd/main.go

wire:
	wire ./internal/di

clean:
	rm -rf bin/
	go clean

docker-up:
	docker run --name postgres-product -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=product_db -p 5432:5432 -d postgres:15

docker-down:
	docker stop postgres-product && docker rm postgres-product

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...
