# Makefile для EncryptKeep Backend

.PHONY: test test-keymanager test-crypto test-codec test-vault test-all build clean test-new

# Цвета для вывода
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

# Запуск всех тестов (новая структура)
test: test-new

# Новая структура тестов
test-new:
	@echo "$(GREEN)Запуск тестов из новой структуры...$(NC)"
	cd backend && make -C tests test

# Тесты для модуля KeyManager (новая структура)
test-keymanager:
	@echo "$(GREEN)Running KeyManager tests...$(NC)"
	cd backend && make -C tests test-keymanager

# Тесты для модуля Crypto (новая структура)
test-crypto:
	@echo "$(GREEN)Running Crypto tests...$(NC)"
	cd backend && make -C tests test-crypto

# Тесты для модуля Codec (новая структура)
test-codec:
	@echo "$(GREEN)Running Codec tests...$(NC)"
	cd backend && make -C tests test-codec

# Тесты для модуля Vault (новая структура)
test-vault:
	@echo "$(GREEN)Running Vault tests...$(NC)"
	cd backend && make -C tests test-vault

# Тесты для модуля Blockchain (новая структура)
test-blockchain:
	@echo "$(GREEN)Running Blockchain tests...$(NC)"
	cd backend && make -C tests test-blockchain

# Запуск всех тестов (старая структура - для совместимости)
test-all:
	@echo "$(YELLOW)Running all tests (old structure)...$(NC)"
	cd backend && go test -v ./internal/...

# Запуск всех тестов (старая структура)
test-old:
	@echo "$(YELLOW)Running all tests (old structure)...$(NC)"
	cd backend && go test -v ./internal/...

# Сборка проекта
build:
	@echo "Building EncryptKeep..."
	cd backend && go build -o encryptkeep.exe ./cmd/encryptkeep

# Очистка
clean:
	@echo "Cleaning up..."
	cd backend && go clean
	rm -f backend/encryptkeep.exe

# Проверка форматирования кода
fmt:
	@echo "Formatting code..."
	cd backend && go fmt ./...

# Проверка линтера
lint:
	@echo "Running linter..."
	cd backend && go vet ./...

# Проверка зависимостей
deps:
	@echo "Checking dependencies..."
	cd backend && go mod tidy
	cd backend && go mod verify

# Установка зависимостей
install:
	@echo "Installing dependencies..."
	cd backend && go mod download

# Показать покрытие тестами
coverage:
	@echo "Running tests with coverage..."
	cd backend && go test -coverprofile=coverage.out ./internal/...
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: backend/coverage.html"

# Быстрый тест (без verbose)
test-quick:
	@echo "Running quick tests..."
	cd backend && go test ./internal/...

# Тест конкретного модуля (использование: make test-module MODULE=keymanager)
test-module:
	@echo "Running tests for module: $(MODULE)"
	cd backend && go test -v ./internal/$(MODULE)

# Помощь
help:
	@echo "Available commands:"
	@echo "  test           - Run all tests"
	@echo "  test-keymanager - Run KeyManager tests"
	@echo "  test-crypto    - Run Crypto tests"
	@echo "  test-codec     - Run Codec tests"
	@echo "  test-vault     - Run Vault tests"
	@echo "  test-all       - Run all tests"
	@echo "  test-quick     - Run tests without verbose output"
	@echo "  test-module    - Run tests for specific module (use MODULE=name)"
	@echo "  build          - Build the project"
	@echo "  clean          - Clean build artifacts"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  deps           - Check dependencies"
	@echo "  install        - Install dependencies"
	@echo "  coverage       - Generate test coverage report"
	@echo "  help           - Show this help"
