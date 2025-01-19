.PHONY: build run test lint swagger clean migrate

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run the application
run:
	go run cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Generate swagger documentation
swagger:
	swag init -g cmd/api/main.go -o docs

# Clean build artifacts
clean:
	rm -rf bin/

# Start database
db-up:
	docker-compose up -d

# Stop database
db-down:
	docker-compose down

# Run database migrations
migrate:
	go run cmd/api/main.go migrate

# Install dependencies
deps:
	go mod download

# Run all pre-commit checks
check: lint test

# Build and run
dev: swagger build run 