# ===================================================================
#  Project configuration
# ===================================================================
PROJECT_NAME    := skinprice
APP_WORKERS     ?= 1
RELEASE_VERSION ?=
RELEASE_DIR     ?= release
RELEASE_ASSETS_DIR ?= $(RELEASE_DIR)/assets
RELEASE_BUILD_DIR ?= ./bin
WAILS ?= wails
WAILS_BUILD_DIR ?= $(PROJECT_NAME)/build/bin

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
        release-clean release-build-linux release-build-windows release-build-binaries release-package-assets release-package-local \
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
	$(GO) install github.com/wailsapp/wails/v2/cmd/wails@latest
	lefthook install


clean: ## Cleanup build artifacts
	rm -rf bin
	$(GO) clean -testcache

release-clean: ## Cleanup generated release packaging artifacts
	rm -rf $(RELEASE_ASSETS_DIR) $(RELEASE_DIR)/bootstrap

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

wails:
	cd $(PROJECT_NAME) && wails dev

release-build-linux: ## Build Linux app + launcher into $(RELEASE_BUILD_DIR)/linux
	mkdir -p "$(RELEASE_BUILD_DIR)/linux"
	cd "$(PROJECT_NAME)" && $(WAILS) build -clean -platform linux/amd64 -tags webkit2_41
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -o "$(RELEASE_BUILD_DIR)/linux/launcher" ./skinprice/cmd/launcher
	cp "$(WAILS_BUILD_DIR)/SkinPrice" "$(RELEASE_BUILD_DIR)/linux/skinprice"
	chmod +x "$(RELEASE_BUILD_DIR)/linux/launcher" "$(RELEASE_BUILD_DIR)/linux/skinprice"

release-build-windows: ## Build Windows app + launcher into $(RELEASE_BUILD_DIR)/windows
	mkdir -p "$(RELEASE_BUILD_DIR)/windows"
	cd "$(PROJECT_NAME)" && CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ $(WAILS) build -clean -platform windows/amd64
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GO) build -o "$(RELEASE_BUILD_DIR)/windows/launcher.exe" ./skinprice/cmd/launcher
	cp "$(WAILS_BUILD_DIR)/SkinPrice.exe" "$(RELEASE_BUILD_DIR)/windows/SkinPrice.exe"

release-build-binaries: ## Build Linux and Windows binaries required for release packaging into $(RELEASE_BUILD_DIR)
	@$(MAKE) release-build-linux RELEASE_BUILD_DIR="$(RELEASE_BUILD_DIR)"
	@$(MAKE) release-build-windows RELEASE_BUILD_DIR="$(RELEASE_BUILD_DIR)"

release-package-assets: ## Build release bootstrap/update packages from $(RELEASE_DIR)/linux and $(RELEASE_DIR)/windows
	@if [ -z "$(RELEASE_VERSION)" ]; then \
		echo "RELEASE_VERSION is required, example: make release-package-assets RELEASE_VERSION=0.2.0"; \
		exit 1; \
	fi
	@set -euo pipefail; \
	mkdir -p "$(RELEASE_ASSETS_DIR)"; \
	tar -C "$(RELEASE_DIR)/linux" -czf "$(RELEASE_ASSETS_DIR)/skinprice-linux-amd64.tar.gz" skinprice; \
	( \
		cd "$(RELEASE_DIR)/windows"; \
		zip -q "$(abspath $(RELEASE_ASSETS_DIR))/skinprice-windows-amd64.zip" SkinPrice.exe; \
	); \
	mkdir -p "$(RELEASE_DIR)/bootstrap/linux/versions/$(RELEASE_VERSION)" "$(RELEASE_DIR)/bootstrap/linux/logs"; \
	cp "$(RELEASE_DIR)/linux/launcher" "$(RELEASE_DIR)/bootstrap/linux/launcher"; \
	cp "$(RELEASE_DIR)/linux/skinprice" "$(RELEASE_DIR)/bootstrap/linux/versions/$(RELEASE_VERSION)/skinprice"; \
	touch "$(RELEASE_DIR)/bootstrap/linux/logs/.keep"; \
	printf '{\n  "version": "%s",\n  "entrypoint": "versions/%s/skinprice",\n  "previous": "",\n  "updated_at": "%s"\n}\n' \
		"$(RELEASE_VERSION)" "$(RELEASE_VERSION)" "$$(date -u +%Y-%m-%dT%H:%M:%SZ)" > "$(RELEASE_DIR)/bootstrap/linux/current.json"; \
	tar -C "$(RELEASE_DIR)/bootstrap/linux" -czf "$(RELEASE_ASSETS_DIR)/skinprice-bootstrap-linux-amd64.tar.gz" .; \
	mkdir -p "$(RELEASE_DIR)/bootstrap/windows/versions/$(RELEASE_VERSION)" "$(RELEASE_DIR)/bootstrap/windows/logs"; \
	cp "$(RELEASE_DIR)/windows/launcher.exe" "$(RELEASE_DIR)/bootstrap/windows/launcher.exe"; \
	cp "$(RELEASE_DIR)/windows/SkinPrice.exe" "$(RELEASE_DIR)/bootstrap/windows/versions/$(RELEASE_VERSION)/SkinPrice.exe"; \
	touch "$(RELEASE_DIR)/bootstrap/windows/logs/.keep"; \
	printf '{\n  "version": "%s",\n  "entrypoint": "versions/%s/SkinPrice.exe",\n  "previous": "",\n  "updated_at": "%s"\n}\n' \
		"$(RELEASE_VERSION)" "$(RELEASE_VERSION)" "$$(date -u +%Y-%m-%dT%H:%M:%SZ)" > "$(RELEASE_DIR)/bootstrap/windows/current.json"; \
	( \
		cd "$(RELEASE_DIR)/bootstrap/windows"; \
		zip -qr "$(abspath $(RELEASE_ASSETS_DIR))/skinprice-bootstrap-windows-amd64.zip" .; \
	); \
	WINDOWS_SHA="$$(sha256sum "$(RELEASE_ASSETS_DIR)/skinprice-windows-amd64.zip" | awk '{print $$1}')"; \
	WINDOWS_SIZE="$$(stat -c %s "$(RELEASE_ASSETS_DIR)/skinprice-windows-amd64.zip")"; \
	LINUX_SHA="$$(sha256sum "$(RELEASE_ASSETS_DIR)/skinprice-linux-amd64.tar.gz" | awk '{print $$1}')"; \
	LINUX_SIZE="$$(stat -c %s "$(RELEASE_ASSETS_DIR)/skinprice-linux-amd64.tar.gz")"; \
	PUBLISHED_AT="$$(date -u +%Y-%m-%dT%H:%M:%SZ)"; \
	printf '{\n  "version": "%s",\n  "channel": "stable",\n  "min_supported_version": "0.1.0",\n  "release_notes": "Release %s",\n  "published_at": "%s",\n  "assets": [\n    {\n      "os": "windows",\n      "arch": "amd64",\n      "filename": "skinprice-windows-amd64.zip",\n      "sha256": "%s",\n      "size": %s,\n      "entrypoint": "SkinPrice.exe"\n    },\n    {\n      "os": "linux",\n      "arch": "amd64",\n      "filename": "skinprice-linux-amd64.tar.gz",\n      "sha256": "%s",\n      "size": %s,\n      "entrypoint": "skinprice"\n    }\n  ]\n}\n' \
		"$(RELEASE_VERSION)" "$(RELEASE_VERSION)" "$$PUBLISHED_AT" "$$WINDOWS_SHA" "$$WINDOWS_SIZE" "$$LINUX_SHA" "$$LINUX_SIZE" > "$(RELEASE_ASSETS_DIR)/update-manifest.json"

release-package-local: ## Build local Windows/Linux binaries into $(RELEASE_BUILD_DIR) and package release assets
	@if [ -z "$(RELEASE_VERSION)" ]; then \
		echo "RELEASE_VERSION is required, example: make release-package-local RELEASE_VERSION=0.2.0"; \
		exit 1; \
	fi
	@$(MAKE) release-build-binaries RELEASE_BUILD_DIR="$(RELEASE_BUILD_DIR)"
	@$(MAKE) release-package-assets RELEASE_VERSION="$(RELEASE_VERSION)" RELEASE_DIR="$(RELEASE_BUILD_DIR)"
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
