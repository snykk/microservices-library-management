# Variables
COMPOSE_FILE = docker-compose.yml
GATEWAY_BINARY=gatewayApp
AUTHOR_BINARY=authorApp
BOOK_BINARY=bookApp
CATEGORY_BINARY=categoryApp
USER_BINARY=userApp
AUTH_BINARY=authApp

# Default target
.PHONY: up down restart logs build clean help

# Build binary project and start the application
up_build: build_gateway build_author build_book build_category build_user, build_auth
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Starting services..."
	docker-compose -f $(COMPOSE_FILE) up --build -d
	@echo "Service started!"

# Stop the application
down:
	@echo "Stopping services..."
	docker-compose -f $(COMPOSE_FILE) down

build_gateway:
	@echo "Building gateway binary..."
	cd ./api-gateway && env GOOS=linux CGO_ENABLED=0 go build -o ${GATEWAY_BINARY} .
	@echo "Done!"

build_author:
	@echo "Building author binary..."
	cd ./author_service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTHOR_BINARY} .
	@echo "Done!"

build_book:
	@echo "Building book binary..."
	cd ./book_service && env GOOS=linux CGO_ENABLED=0 go build -o ${BOOK_BINARY} .
	@echo "Done!"

build_category:
	@echo "Building category binary..."
	cd ./category_service && env GOOS=linux CGO_ENABLED=0 go build -o ${CATEGORY_BINARY} .
	@echo "Done!"

build_user:
	@echo "Building user binary..."
	cd ./user_service && env GOOS=linux CGO_ENABLED=0 go build -o ${USER_BINARY} .
	@echo "Done!"

build_auth:
	@echo "Building auath binary..."
	cd ./auath_service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} .
	@echo "Done!"

# Restart the application
restart: down up

# View logs
logs:
	@echo "Viewing logs..."
	docker-compose -f $(COMPOSE_FILE) logs -f

# Build images
build:
	@echo "Building services..."
	docker-compose -f $(COMPOSE_FILE) build

# Clean up unused images and volumes
clean:
	@echo "Cleaning up unused resources..."
	docker-compose -f $(COMPOSE_FILE) down --volumes --remove-orphans

# Display help
help:
	@echo "Usage:"
	@echo "  make up       Start all services in detached mode"
	@echo "  make down     Stop all services and remove containers"
	@echo "  make restart  Restart all services"
	@echo "  make logs     View logs for all services"
	@echo "  make build    Build images for all services"
	@echo "  make clean    Remove containers, volumes, and orphans"
