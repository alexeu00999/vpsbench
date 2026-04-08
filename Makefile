# --- Makefile для VPSBENCH (vpsbench) ---
# Использование: make [target]

SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

# --- Проект ---
PROJECT   ?= vpsbench
GO        ?= go
GOFLAGS   ?=
LDFLAGS   ?= -s -w -X 'github.com/user/vpsbench/internal/sysinfo.Version=$(VERSION)' -X 'github.com/user/vpsbench/internal/sysinfo.Commit=$(COMMIT)' -X 'github.com/user/vpsbench/internal/sysinfo.BuildTime=$(BUILD_TIME)'

# --- Git ---
VERSION   ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT    ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# --- Сборка ---
BIN_DIR    := bin
MAIN_PKG   ?= ./cmd/vpsbench
BINARY     := $(BIN_DIR)/$(PROJECT)

# --- Docker ---
DOCKER_IMAGE ?= vpsbench/vpsbench
DOCKER_TAG   ?= $(VERSION)

# --- Инструменты ---
GOLANGCI_LINT ?= golangci-lint
GOTEST        ?= $(GO) test
GOTESTFLAGS   ?= -race -count=1

# --- Платформы для кросс-компиляции ---
PLATFORMS ?= linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

# ============================================================================
.DEFAULT_GOAL := help

##@ Разработка

.PHONY: build
build: ## Собрать бинарник
	$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BINARY) $(MAIN_PKG)

.PHONY: build-linux-amd64
build-linux-amd64: ## Собрать бинарник для Linux x86_64
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(PROJECT)-linux-amd64 $(MAIN_PKG)

.PHONY: build-linux-arm64
build-linux-arm64: ## Собрать бинарник для Linux ARM64
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(PROJECT)-linux-arm64 $(MAIN_PKG)

.PHONY: build-darwin-amd64
build-darwin-amd64: ## Собрать бинарник для macOS Intel
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(PROJECT)-darwin-amd64 $(MAIN_PKG)

.PHONY: build-darwin-arm64
build-darwin-arm64: ## Собрать бинарник для macOS Apple Silicon
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(PROJECT)-darwin-arm64 $(MAIN_PKG)

.PHONY: run
run: build ## Собрать и запустить
	$(BINARY)

.PHONY: dev
dev: ## Запуск с hot reload (требуется air)
	air

.PHONY: install
install: build ## Установить в $GOPATH/bin
	cp $(BINARY) $(shell $(GO) env GOPATH)/bin/$(PROJECT)

##@ Тестирование

.PHONY: test
test: ## Запустить тесты
	$(GOTEST) $(GOTESTFLAGS) ./...

.PHONY: test-cover
test-cover: ## Тесты с отчётом покрытия
	$(GOTEST) $(GOTESTFLAGS) -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Отчёт покрытия: coverage.html"

.PHONY: bench
bench: ## Запустить Go-бенчмарки
	$(GO) test -bench=. -benchmem ./...

##@ Качество кода

.PHONY: lint
lint: ## Запустить линтеры
	$(GOLANGCI_LINT) run ./...

.PHONY: fmt
fmt: ## Отформатировать код
	$(GO) fmt ./...
	goimports -w .

.PHONY: vet
vet: ## Запустить go vet
	$(GO) vet ./...

.PHONY: tidy
tidy: ## Привести в порядок go.mod
	$(GO) mod tidy
	$(GO) mod verify

##@ Docker

.PHONY: docker-build
docker-build: ## Собрать Docker-образ
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		-t $(DOCKER_IMAGE):latest \
		.

.PHONY: docker-run
docker-run: ## Запустить в Docker
	docker run --rm -it $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: docker-push
docker-push: ## Запушить Docker-образ
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_IMAGE):latest

##@ Релиз

.PHONY: build-all
build-all: clean ## Кросс-компиляция для всех платформ
	@mkdir -p $(BIN_DIR)
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		output=$(BIN_DIR)/$(PROJECT)-$${os}-$${arch}; \
		echo "Сборка: $${os}/$${arch} → $${output}"; \
		CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $${output} $(MAIN_PKG) || exit $$?; \
	done

.PHONY: checksum
checksum: build-all ## Сгенерировать контрольные суммы (sha256)
	@echo "Генерация контрольных сумм..."
	cd $(BIN_DIR) && sha256sum $(PROJECT)-* > checksums.txt
	@cat $(BIN_DIR)/checksums.txt

##@ CI

.PHONY: ci
ci: lint test build ## Полный CI пайплайн

##@ Очистка

.PHONY: clean
clean: ## Удалить артефакты сборки
	rm -rf $(BIN_DIR) coverage.out coverage.html

##@ Помощь

.PHONY: help
help: ## Показать эту справку
	@awk 'BEGIN {FS = ":.*##"; printf "Использование:\n  make \033[36m<target>\033[0m\n"} \
		/^[a-zA-Z0-9_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2} \
		/^##@/ {printf "\n\033[1m%s\033[0m\n", substr($$0, 5)}' $(MAKEFILE_LIST)
