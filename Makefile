# Makefile v26.03.22
# Версии инструментов
BUF_VERSION := latest
PROTOC_GEN_GO_VERSION := latest
PROTOC_GEN_GO_GRPC_VERSION := latest
PROTOC_GEN_VALIDATE_VERSION := latest
PROTOC_GRPC_GATEWAY_VERSION := latest
PROTOC_GEN_OPENAPI_V2_VERSION := latest
GOLANGCI_LINT_VERSION := latest
MOCKERY_VERSION := latest
SWAG_VERSION := latest
SCC_VERSION := latest

# Директории
ROOT_DIR := $(shell pwd)
BIN_DIR ?= $(or $(GOBIN),$(GOPATH)/bin,$(HOME)/go/bin)
PROTO_DIR := $(ROOT_DIR)/proto

# Экспортируем GOBIN, чтобы go install клал бинарники в BIN_DIR
export GOBIN := $(BIN_DIR)
# Добавляем GOBIN в PATH, чтобы buf находил protoc-плагины
export PATH := $(BIN_DIR):$(PATH)

# Бинарные файлы
BUF := $(BIN_DIR)/buf
PROTOC_GEN_GO := $(BIN_DIR)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(BIN_DIR)/protoc-gen-go-grpc
PROTOC_GEN_VALIDATE := $(BIN_DIR)/protoc-gen-validate
PROTOC_GEN_GRPC_GATEWAY := $(BIN_DIR)/protoc-gen-grpc-gateway
PROTOC_GEN_OPENAPI_V2 := $(BIN_DIR)/protoc-gen-openapiv2
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint
MOCKERY := $(BIN_DIR)/mockery

# Цели, которые не создают файлы
.PHONY: all install-buf install-golangci-lint install-mockery proto-install-plugins proto-lint proto-update-deps proto-gen mock-gen lint test help

# Цель по умолчанию
all: proto-gen

# Справка
help:
	@echo "Доступные команды:"
	@echo "  make install-buf           - Устанавливает Buf в \$$GOBIN"
	@echo "  make install-golangci-lint - Устанавливает golangci-lint в \$$GOBIN"
	@echo "  make install-mockery       - Устанавливает mockery в \$$GOBIN"
	@echo "  make proto-install-plugins - Устанавливает protoc плагины в \$$GOBIN"
	@echo "  make proto-lint            - Проверка .proto-файлов на соответствие стилю"
	@echo "  make proto-update-deps     - Обновляет зависимости protobuf из удаленных репозиториев"
	@echo "  make proto-gen             - Генерация Go-кода из .proto"
	@echo "  make mock-gen              - Генерация mock-объектов через mockery"
	@echo "  make lint                  - Go Linter"
	@echo "  make test                  - Запуск тестов с race detection и coverage"
	@echo ""
	@echo "  GOBIN = $(BIN_DIR)"

# Утилита: устанавливает инструмент, если его нет в $GOBIN
# $(1) — имя бинарника, $(2) — go module path, $(3) — версия
define install-tool
	@mkdir -p $(BIN_DIR)
	@if [ ! -f $(BIN_DIR)/$(1) ]; then \
		echo "Устанавливаем $(1)@$(3) в $(BIN_DIR)..."; \
		go install $(2)@$(3); \
	else \
		echo "$(1) уже установлен"; \
	fi
endef

# Установка swag-go
install-swag:
	$(call install-tool,scc,github.com/swaggo/swag/cmd/swag,$(SWAG_VERSION))

# Установка scc
install-scc:
	$(call install-tool,scc,github.com/boyter/scc/v3,$(SCC_VERSION))

# Установка Buf
install-buf:
	$(call install-tool,buf,github.com/bufbuild/buf/cmd/buf,$(BUF_VERSION))

# Установка golangci-lint
install-golangci-lint:
	$(call install-tool,golangci-lint,github.com/golangci/golangci-lint/v2/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

# Установка mockery
install-mockery:
	$(call install-tool,mockery,github.com/vektra/mockery/v3,$(MOCKERY_VERSION))

# Установка protoc плагинов
install-protoc-gen-go:
	$(call install-tool,protoc-gen-go,google.golang.org/protobuf/cmd/protoc-gen-go,$(PROTOC_GEN_GO_VERSION))

install-protoc-gen-go-grpc:
	$(call install-tool,protoc-gen-go-grpc,google.golang.org/grpc/cmd/protoc-gen-go-grpc,$(PROTOC_GEN_GO_GRPC_VERSION))

install-protoc-gen-validate:
	$(call install-tool,protoc-gen-validate,github.com/envoyproxy/protoc-gen-validate,$(PROTOC_GEN_VALIDATE_VERSION))

install-protoc-gen-grpc-gateway:
	$(call install-tool,protoc-gen-grpc-gateway,github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway,$(PROTOC_GRPC_GATEWAY_VERSION))

install-protoc-gen-openapiv2:
	$(call install-tool,protoc-gen-openapiv2,github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2,$(PROTOC_GEN_OPENAPI_V2_VERSION))

# Установка всех protoc плагинов
proto-install-plugins: install-protoc-gen-go install-protoc-gen-go-grpc install-protoc-gen-validate install-protoc-gen-grpc-gateway install-protoc-gen-openapiv2
	@echo "Все protoc плагины установлены"

# Проверка .proto-файлов на соответствие стилю
proto-lint: install-buf proto-install-plugins
	@echo "Проверка .proto-файлов на соответствие стилю..."
	@cd $(PROTO_DIR) && $(BUF) lint

# Обновление зависимостей protobuf
proto-update-deps: install-buf
	@echo "Обновляем зависимости buf..."
	@cd $(PROTO_DIR) && $(BUF) dep update

# Обновление зависимостей golang
project-update:
	go get all
	go mod tidy

# Генерация Go-кода из .proto
proto-gen: install-buf proto-install-plugins proto-lint
	@echo "Генерация Go-кода из .proto..."
	@cd $(PROTO_DIR) && $(BUF) generate

# Генерация mock-объектов через mockery
mock-gen: install-mockery
	@echo "Генерация mock-объектов через mockery..."
	@$(MOCKERY) --config ./env/mockery/mockery.yml
	@echo "Mock-объекты успешно сгенерированы"

# Go Linter
lint: install-golangci-lint
	@echo "Запуск Go Linter..."
	@$(GOLANGCI_LINT) run --config ./env/lint/golangci.yml ./...

# Go SCC, analog is below
# https://github.com/XAMPPRocky/tokei
# https://github.com/AlDanial/cloc
scc: install-scc
	@echo "Запуск Sloc Cloc and Code..."
	@scc .

swag: install-swag
	@swag fmt
	@swag init -g ./internal/api/rest/handlers.go


# Запуск всех видов тестов
test:
	@echo "Запуск тестов с race detection и coverage..."
	@go test -v -race -cover -coverprofile=coverage.out ./...
	@echo ""
	@echo "Покрытие кода:"
	@go tool cover -func=coverage.out

build: proto-gen lint test
	@echo "Компиляция..."
	@go build ./cmd/app/main.go
