.PHONY: build run test clean migrate docker-up docker-down

APP_NAME := genvid-backend
VERSION := 1.0.0
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_VERSION := $(shell go version | awk '{print $$3}')

build:
	go build -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" -o bin/$(APP_NAME) ./cmd/server

run:
	go run ./cmd/server

dev:
	go run ./cmd/server

test:
	go test -v -race -coverprofile=coverage.out ./...

test-coverage:
	go tool cover -html=coverage.out

clean:
	rm -rf bin/
	rm -f coverage.out

deps:
	go mod download
	go mod tidy

lint:
	golangci-lint run

migrate-up:
	@echo "Run migrations with: psql -d genvid -f migrations/001_initial_schema.sql"

migrate-down:
	@echo "Rollback not implemented"

docker-build:
	docker build -t $(APP_NAME):$(VERSION) .

docker-run:
	docker run -p 8080:8080 --env-file .env $(APP_NAME):$(VERSION)

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

db-connect:
	psql $(DATABASE_URL)

format:
	go fmt ./...

vet:
	go vet ./...

all: deps format vet test build
