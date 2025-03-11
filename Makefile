.PHONY: all build run test clean deps lint migrate dev docker-up docker-down

all: build

build:
	go build -o bin/book-library ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v ./...

clean:
	rm -rf bin/

deps:
	go mod download

lint:
	golangci-lint run ./...

migrate:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/book_library?sslmode=disable" up

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

migrate-down:
	@read -p "Enter number of down migration: " n; \
	migrate -database "postgres://postgres:postgres@localhost:5432/book_library?sslmode=disable" -path ./migrations down $$n

dev:
	go run ./cmd/server

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-restart:
	docker-compose down
	docker-compose up -d --build

mocks:
	mockgen -source=internal/repository/interfaces/repository_interfaces.go -destination=internal/mocks/repository_mocks.go -package=mocks

help:
	@echo "Available targets:"
	@echo "  all          - Build the application (default)"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  lint         - Run linter"
	@echo "  migrate      - Run migrations"
	@echo "  migrate-create - Create a new migration"
	@echo "  dev          - Start development environment"
	@echo "  docker-up    - Start Docker containers"
	@echo "  docker-down  - Stop Docker containers"
	@echo "  docker-restart - Rebuild and restart Docker containers"
	@echo "  mocks        - Generate mocks"
	@echo "  help         - Show this help"