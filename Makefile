.PHONY: help pull up down restart logs logs-gateway logs-auth logs-content logs-media logs-chat ps health clean rebuild update-env

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Show this help message
	@echo "$(BLUE)Writeful Backend Deployment Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

pull: ## Pull latest images from Docker Hub
	@echo "$(BLUE)Pulling latest images from Docker Hub...$(NC)"
	docker-compose pull

up: ## Start all services
	@echo "$(BLUE)Starting all services...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)All services started!$(NC)"
	@echo "$(YELLOW)Gateway available at: http://localhost:8080$(NC)"
	@make ps

down: ## Stop all services
	@echo "$(BLUE)Stopping all services...$(NC)"
	docker-compose down
	@echo "$(GREEN)All services stopped!$(NC)"

restart: ## Restart all services
	@echo "$(BLUE)Restarting all services...$(NC)"
	docker-compose restart
	@echo "$(GREEN)All services restarted!$(NC)"

restart-gateway: ## Restart and recreate gateway service
	@echo "$(BLUE)Recreating gateway service container...$(NC)"
	docker-compose up -d --force-recreate gateway-service
	@echo "$(GREEN)Gateway service recreated and restarted!$(NC)"

restart-auth: ## Restart and recreate auth service
	@echo "$(BLUE)Recreating auth service container...$(NC)"
	docker-compose up -d --force-recreate auth-service
	@echo "$(GREEN)Auth service recreated and restarted!$(NC)"

restart-content: ## Restart and recreate content service
	@echo "$(BLUE)Recreating content service container...$(NC)"
	docker-compose up -d --force-recreate content-service
	@echo "$(GREEN)Content service recreated and restarted!$(NC)"

restart-media: ## Restart and recreate media service
	@echo "$(BLUE)Recreating media service container...$(NC)"
	docker-compose up -d --force-recreate media-service
	@echo "$(GREEN)Media service recreated and restarted!$(NC)"

restart-chat: ## Restart and recreate chat service
	@echo "$(BLUE)Recreating chat service container...$(NC)"
	docker-compose up -d --force-recreate chat-service
	@echo "$(GREEN)Chat service recreated and restarted!$(NC)"

logs: ## Show logs for all services
	docker-compose logs -f

logs-gateway: ## Show logs for gateway service
	docker-compose logs -f gateway-service

logs-auth: ## Show logs for auth service
	docker-compose logs -f auth-service

logs-content: ## Show logs for content service
	docker-compose logs -f content-service

logs-media: ## Show logs for media service
	docker-compose logs -f media-service

logs-chat: ## Show logs for chat service
	docker-compose logs -f chat-service

ps: ## Show status of all services
	@echo "$(BLUE)Service Status:$(NC)"
	@docker-compose ps

health: ## Check health of gateway
	@echo "$(BLUE)Checking gateway health...$(NC)"
	@curl -s http://localhost:8080/health | jq . || echo "$(RED)Gateway is not responding$(NC)"

clean: ## Stop and remove all containers, networks
	@echo "$(YELLOW)Stopping and removing all containers...$(NC)"
	docker-compose down
	@echo "$(GREEN)Cleanup complete!$(NC)"

rebuild: ## Pull latest images and recreate containers
	@echo "$(BLUE)Pulling latest images...$(NC)"
	docker-compose pull
	@echo "$(BLUE)Recreating containers...$(NC)"
	docker-compose up -d --force-recreate
	@echo "$(GREEN)Rebuild complete!$(NC)"
	@make ps

deploy: pull up ## Deploy all services (pull + up)
	@echo "$(GREEN)Deployment complete!$(NC)"
	@echo "$(YELLOW)Gateway: http://localhost:8080$(NC)"
	@echo "$(YELLOW)Auth Service: http://localhost:8004$(NC)"
	@echo "$(YELLOW)Content Service: http://localhost:8003$(NC)"
	@echo "$(YELLOW)Media Service: http://localhost:8005$(NC)"
	@echo "$(YELLOW)Chat Service: http://localhost:8006$(NC)"

build-all: ## Build all local service images using docker-compose
	@echo "$(BLUE)Building all local service images...$(NC)"
	docker-compose build

push-all: build-all ## Build and push all service images to Docker Hub
	@echo "$(BLUE)Pushing all service images to Docker Hub...$(NC)"
	docker-compose push
	@echo "$(GREEN)All images pushed successfully!$(NC)"

tidy-all: ## Run go mod tidy on all microservices to update go.sum
	@echo "$(BLUE)Tidying auth-service...$(NC)"
	cd auth-service && go mod tidy
	@echo "$(BLUE)Tidying chat-service...$(NC)"
	cd chat-service && go mod tidy
	@echo "$(BLUE)Tidying content-service...$(NC)"
	cd content-service && go mod tidy
	@echo "$(GREEN)All go.sum files updated!$(NC)"

check-db: ## Check database connection
	@echo "$(BLUE)Checking database connection...$(NC)"
	@psql -h localhost -p 5438 -U postgres -d playground -c "SELECT 'Database is accessible!' as status;" || echo "$(RED)Cannot connect to database$(NC)"

check-ports: ## Check if ports are available
	@echo "$(BLUE)Checking ports...$(NC)"
	@lsof -i :8080 || echo "$(GREEN)Port 8080 is available$(NC)"
	@lsof -i :8004 || echo "$(GREEN)Port 8004 is available$(NC)"
	@lsof -i :8003 || echo "$(GREEN)Port 8003 is available$(NC)"
	@lsof -i :8005 || echo "$(GREEN)Port 8005 is available$(NC)"
	@lsof -i :8006 || echo "$(GREEN)Port 8006 is available$(NC)"

test-endpoints: ## Test all service endpoints
	@echo "$(BLUE)Testing Gateway Health...$(NC)"
	@curl -s http://localhost:8080/health | jq . || echo "$(RED)Gateway health check failed$(NC)"
	@echo ""
	@echo "$(BLUE)All tests complete!$(NC)"

update-env: ## Update Cloudflare Tunnel URLs to Render FE and redeploy
	@./update-render-env.sh

install-proto-tools: ## Install protoc and Go gRPC generators
	@echo "Installing Go plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Installing protoc compiler..."
	mkdir -p bin_temp
	curl -L https://github.com/protocolbuffers/protobuf/releases/download/v27.0/protoc-27.0-osx-aarch_64.zip -o bin_temp/protoc.zip
	unzip -o bin_temp/protoc.zip -d bin_temp
	mkdir -p /Users/haopham/go/bin
	mkdir -p /Users/haopham/go/include
	mv bin_temp/bin/protoc /Users/haopham/go/bin/protoc
	cp -R bin_temp/include/ /Users/haopham/go/include/
	rm -rf bin_temp
	@echo "Tools installed successfully!"
	@make gen-proto

gen-proto: ## Generate Go code from proto files
	@echo "Generating protobuf files..."
	mkdir -p auth-service/pkg/pb/auth
	mkdir -p chat-service/pkg/pb/auth
	mkdir -p content-service/pkg/pb/auth
	/Users/haopham/go/bin/protoc --go_out=auth-service/pkg/pb/auth --go_opt=paths=source_relative \
	       --go-grpc_out=auth-service/pkg/pb/auth --go-grpc_opt=paths=source_relative \
	       -I proto/auth -I /Users/haopham/go/include proto/auth/auth.proto
	/Users/haopham/go/bin/protoc --go_out=chat-service/pkg/pb/auth --go_opt=paths=source_relative \
	       --go-grpc_out=chat-service/pkg/pb/auth --go-grpc_opt=paths=source_relative \
	       -I proto/auth -I /Users/haopham/go/include proto/auth/auth.proto
	/Users/haopham/go/bin/protoc --go_out=content-service/pkg/pb/auth --go_opt=paths=source_relative \
	       --go-grpc_out=content-service/pkg/pb/auth --go-grpc_opt=paths=source_relative \
	       -I proto/auth -I /Users/haopham/go/include proto/auth/auth.proto
	@echo "Protobuf files generated successfully!"



