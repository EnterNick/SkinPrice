# ===================================================================
#  Project configuration
# ===================================================================
PROJECT_NAME    := skinprice
APP_WORKERS     ?= 1

DOCKER_FILENAME ?= docker-compose.dev.yaml

# Go
GO              ?= go
GOFLAGS         ?=
PKG             ?= ./...
CMD_DIR         ?= ./$(PROJECT_NAME)/cmd/
MAIN_PKG        ?= $(CMD_DIR)main/main.go
MIGRATION_NAME ?= init

# Lint/Format tools (must exist in PATH)
GOLANGCI_LINT   ?= golangci-lint
GOIMPORTS       ?= goimports
GOOSE           ?= goose

# Docker image
HARBOR_USERNAME ?=
HARBOR_PASSWORD ?=
HARBOR_REGISTRY ?=
IMAGE_NAME      ?= $(PROJECT_NAME)/backend/$(PROJECT_NAME)
TAG             ?= $(shell git rev-parse --short=8 HEAD)

# Env encryption
ENCRYPTED_FILE   ?= .env.enc
DECRYPTED_FILE   ?= .env.prod
DECRYPTED_SECRET ?=

# ===================================================================
#  Misc
# ===================================================================
SHELL := /usr/bin/env bash

.PHONY: help clean local local_down \
        deps tools \
        fmt format lint lint-ci \
        test test-ci \
        build run dev \
        local-create-migrations local-apply-migrations local-delete-migrations local-recreate-migrations \
        docker-apply-migrations \
        encrypt_env decrypt_env build-ci

help: ## Show this help
	@echo "Available commands:"
	@grep -E "^[a-zA-Z0-9_.-]+:.*##" $(MAKEFILE_LIST) | sed -E "s/:.*##/  /"

# ===================================================================
#  Tooling
# ===================================================================

deps: ## Download Go deps
	$(GO) mod download

tools: ## Install dev tools (goose, goimports, golangci-lint)
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	$(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	$(GO) install github.com/evilmartians/lefthook/v2@latest
	lefthook install


clean: ## Cleanup build artifacts
	rm -rf bin
	$(GO) clean -testcache

# ===================================================================
#  Local stack (Docker)
# ===================================================================

local: ## Start local stack (build & recreate)
	docker compose -f $(DOCKER_FILENAME) up --build --force-recreate --remove-orphans --renew-anon-volumes

local_down: ## Stop local stack and remove volumes
	docker compose -f $(DOCKER_FILENAME) down -v

local-apply-migrations:
	$(GO) run $(CMD_DIR)migrations/apply/main.go

local-create-migrations:
	cd $(PROJECT_NAME) && go run ./cmd/migrations/create/main.go -name $(MIGRATION_NAME)

local-delete-migrations:
	find $(PROJECT_NAME)/internal/adapters/database/ent/migrate/migrations -type f -delete

local-generate-schema:
	cd "$(PROJECT_NAME)" && go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/execquery --target ./internal/adapters/database/ent ./internal/adapters/database/ent/schema

NAME ?=
local-add-table:
	cd $(PROJECT_NAME) && $(GO) run -mod=mod entgo.io/ent/cmd/ent new --target ./internal/adapters/database/ent/schema $(NAME)

local-recreate-migrations: local-delete-migrations local-create-migrations local-apply-migrations
# ===================================================================
#  Quality
# ===================================================================

fmt: ## Format code (gofmt)
	$(GO) fmt $(PKG)

format: ## Format code (goimports)
	$(GOIMPORTS) -w $$(find . -type f -name '*.go' -not -path './vendor/*')

lint: fmt format ## golangci-lint (local)
	$(GOLANGCI_LINT) run --fix --timeout=5m

lint-ci: ## golangci-lint (CI)
	$(GOLANGCI_LINT) run --timeout=5m --out-format=colored-line-number

test: ## Run tests (verbose)
	$(GO) test $(PKG) -v

test-ci: ## Run tests with junit (requires go-junit-report)
	@command -v go-junit-report >/dev/null 2>&1 || { \
	  echo "go-junit-report not found. Installing..."; \
	  $(GO) install github.com/jstemmer/go-junit-report/v2@latest; \
	}
	$(GO) test $(PKG) -v 2>&1 | go-junit-report -set-exit-code > junit.xml

# ===================================================================
#  App
# ===================================================================

build: ## Build binary into ./bin
	mkdir -p bin
	$(GO) build $(GOFLAGS) -o bin/$(PROJECT_NAME) $(MAIN_PKG)

run: ## Run app
	$(GO) run $(GOFLAGS) $(MAIN_PKG)

dev: ## Run with file-watcher (requires air)
	@command -v air >/dev/null 2>&1 || { \
	  echo "air not found. Installing..."; \
	  $(GO) install github.com/air-verse/air@latest; \
	}
	air

# ===================================================================
#  Docker build/push
# ===================================================================

build-docker: ## Build and push Docker image
	@echo -n "$(HARBOR_PASSWORD)" | docker login --username $(HARBOR_USERNAME) --password-stdin $(HARBOR_REGISTRY) || { echo "Docker login failed"; exit 1; }
	@docker build -t $(HARBOR_REGISTRY)/$(IMAGE_NAME):$(TAG) -f Dockerfile . || { echo "Docker build failed"; exit 1; }
	@docker push $(HARBOR_REGISTRY)/$(IMAGE_NAME):$(TAG)
	@echo "Pushed: $(HARBOR_REGISTRY)/$(IMAGE_NAME):$(TAG)"
	@docker logout $(HARBOR_REGISTRY)

# ===================================================================
#  Env encryption
# ===================================================================

encrypt_env: ## Encrypt environment variables file
	openssl enc -aes-256-cbc -salt -in $(DECRYPTED_FILE) -out $(ENCRYPTED_FILE) -k $(DECRYPTED_SECRET)
	@echo "$(DECRYPTED_FILE) encrypted as $(ENCRYPTED_FILE)"

decrypt_env: ## Decrypt environment variables file
	openssl enc -aes-256-cbc -d -in $(ENCRYPTED_FILE) -out $(DECRYPTED_FILE) -k $(DECRYPTED_SECRET)
	@echo "$(ENCRYPTED_FILE) decrypted to $(DECRYPTED_FILE)"

build-ci: ## decrypt_env + build-docker
	@$(MAKE) decrypt_env
	@$(MAKE) build-docker
