# Variables
DOCKER_REGISTRY ?= harbor.remystorage.ru
DOCKER_ORG ?= chaika
IMAGE_NAME ?= chaika-goods
VERSION ?= $(shell cat VERSION)
DOCKER_IMAGE = $(DOCKER_REGISTRY)/$(DOCKER_ORG)/$(IMAGE_NAME)

# Docker build flags
DOCKER_BUILD_FLAGS ?= --no-cache

.PHONY: build push all clean help

help: ## Display this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk -F ':|##' '/^[^\t].+?:.*?##/ { printf "  %-20s %s\n", $$1, $$NF }' $(MAKEFILE_LIST)

build: ## Build the Docker image
	@echo "Building $(DOCKER_IMAGE):$(VERSION)"
	docker build $(DOCKER_BUILD_FLAGS) \
		--build-arg VERSION=$(VERSION) \
		-t $(DOCKER_IMAGE):$(VERSION) \
		-t $(DOCKER_IMAGE):latest \
		.

push: ## Push the Docker image to registry
	@echo "Pushing $(DOCKER_IMAGE):$(VERSION)"
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_IMAGE):latest

all: build push ## Build and push the Docker image

clean: ## Remove local Docker images
	@echo "Removing $(DOCKER_IMAGE):$(VERSION)"
	docker rmi $(DOCKER_IMAGE):$(VERSION) || true
	docker rmi $(DOCKER_IMAGE):latest || true

# Development commands
run: build ## Build and run the container locally
	docker run -p 8080:8080 -p 50051:50051 $(DOCKER_IMAGE):$(VERSION)

test: ## Run tests
	go test ./...

lint: ## Run linter
	golangci-lint run

# Version management
version: ## Show current version
	@echo $(VERSION)

bump-version: ## Bump version (usage: make bump-version VERSION=1.0.2)
	@echo $(VERSION) > VERSION
	@echo "Version bumped to $(VERSION)"

swag: ## Generate Swagger documentation
	swag init -g cmd/main.go --parseInternal
	@echo "Swagger documentation generated in docs/swagger"