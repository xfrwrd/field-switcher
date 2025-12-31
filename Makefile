.PHONY: help test lint check install-tools clean build run

GOLANGCI_LINT_VERSION := v1.55.2
GOLANGCI_LINT_BIN := $(shell go env GOPATH)/bin/golangci-lint

help:
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

test:
	@echo "Запуск тестов..."
	@go test -v -race -coverprofile=coverage.out -coverpkg=./... ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Покрытие кода сохранено в coverage.html"

test-short:
	@echo "Запуск тестов (short)..."
	@go test -v ./...

lint:
	@echo "Запуск линтера..."
	@if [ ! -f "$(GOLANGCI_LINT_BIN)" ]; then \
		echo "golangci-lint не установлен. Запустите: make install-tools"; \
		exit 1; \
	fi
	@$(GOLANGCI_LINT_BIN) run ./...

check: test lint
	@echo "Все проверки пройдены!"

install-tools:
	@echo "Установка golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)
	@echo "golangci-lint установлен: $(GOLANGCI_LINT_BIN)"

build:
	@echo "Сборка приложения..."
	@go build -o bin/field-switcher ./cmd/main.go
	@echo "Приложение собрано: bin/field-switcher"

run:
	@go run ./cmd/main.go

clean:
	@echo "Очистка артефактов..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Очистка завершена"

fmt:
	@echo "Форматирование кода..."
	@go fmt ./...
	@echo "Форматирование завершено"

vet: 
	@echo "Запуск go vet..."
	@go vet ./...
	@echo "go vet завершен"

