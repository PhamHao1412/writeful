.PHONY: help build run stop restart clean logs test migrate-up migrate-down
AUTHOR:= haopham1709
APP_NAME := chat-service
PORT := 8006
DOCKER_IMAGE := $(AUTHOR)/$(APP_NAME):latest
CONTAINER_NAME := $(APP_NAME)

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

push: build ## Build and push Docker image to registry
	@echo "Pushing Docker image to registry..."
	docker push $(DOCKER_IMAGE)

deploy: push ## Build, push and deploy locally via docker-compose
	@echo "Recreating container in docker-compose..."
	docker compose -f ../docker-compose.yml up -d --force-recreate $(CONTAINER_NAME)

run: stop ## Stop and run container
	@echo "Starting $(CONTAINER_NAME)..."
	docker run -d \
		--name $(CONTAINER_NAME) \
		-p $(PORT):$(PORT) \
		$(DOCKER_IMAGE)
	@echo "$(CONTAINER_NAME) is running on port $(PORT)"

stop: ## Stop and remove container
	@echo "Stopping $(CONTAINER_NAME)..."
	@docker stop $(CONTAINER_NAME) 2>/dev/null || true
	@docker rm -f $(CONTAINER_NAME) 2>/dev/null || true

restart: stop run ## Restart container

clean: stop ## Remove container and image
	@echo "Removing Docker image..."
	@docker rmi $(DOCKER_IMAGE) 2>/dev/null || true

logs: ## Show container logs
	docker logs -f $(CONTAINER_NAME)

dev: ## Run locally without Docker
	go run cmd/serverd/main.go

test: ## Run tests
	go test -v ./...

migrate-up: ## Run database migrations
	goose -dir data/migrations postgres "$(shell grep DB.URL .env | cut -d'=' -f2)" up

migrate-down: ## Rollback database migrations
	goose -dir data/migrations postgres "$(shell grep DB.URL .env | cut -d'=' -f2)" down

build-run: build run ## Build and run in one command
