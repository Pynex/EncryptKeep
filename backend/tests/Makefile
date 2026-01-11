# Makefile для тестов EncryptKeep Backend

.PHONY: help test test-unit test-integration test-e2e test-cover test-race test-verbose clean

# Цвета для вывода
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

# Переменные
TEST_DIR=./tests
UNIT_DIR=$(TEST_DIR)/unit
INTEGRATION_DIR=$(TEST_DIR)/integration
E2E_DIR=$(TEST_DIR)/e2e
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

help: ## Показать справку
	@echo "$(GREEN)EncryptKeep Backend Tests$(NC)"
	@echo ""
	@echo "$(YELLOW)Доступные команды:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

test: ## Запустить все тесты
	@echo "$(GREEN)Запуск всех тестов...$(NC)"
	@go test ./...

test-unit: ## Запустить модульные тесты
	@echo "$(GREEN)Запуск модульных тестов...$(NC)"
	@go test ./unit/...

test-integration: ## Запустить интеграционные тесты
	@echo "$(GREEN)Запуск интеграционных тестов...$(NC)"
	@go test ./integration/...

test-e2e: ## Запустить E2E тесты
	@echo "$(GREEN)Запуск E2E тестов...$(NC)"
	@go test ./e2e/...

test-verbose: ## Запустить тесты с подробным выводом
	@echo "$(GREEN)Запуск тестов с подробным выводом...$(NC)"
	@go test -v ./...

test-race: ## Запустить тесты с race detection
	@echo "$(GREEN)Запуск тестов с race detection...$(NC)"
	@go test -race ./...

test-cover: ## Запустить тесты с покрытием кода
	@echo "$(GREEN)Запуск тестов с покрытием кода...$(NC)"
	@go test -coverprofile=$(COVERAGE_FILE) ./...
	@echo "$(GREEN)Покрытие кода:$(NC)"
	@go tool cover -func=$(COVERAGE_FILE)

test-cover-html: test-cover ## Создать HTML отчет о покрытии
	@echo "$(GREEN)Создание HTML отчета о покрытии...$(NC)"
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "$(GREEN)HTML отчет создан: $(COVERAGE_HTML)$(NC)"

test-benchmark: ## Запустить бенчмарки
	@echo "$(GREEN)Запуск бенчмарков...$(NC)"
	@go test -bench=. ./...

test-short: ## Запустить короткие тесты
	@echo "$(GREEN)Запуск коротких тестов...$(NC)"
	@go test -short ./...

test-timeout: ## Запустить тесты с таймаутом
	@echo "$(GREEN)Запуск тестов с таймаутом 30s...$(NC)"
	@go test -timeout=30s ./...

test-parallel: ## Запустить тесты параллельно
	@echo "$(GREEN)Запуск тестов параллельно...$(NC)"
	@go test -parallel=4 ./...

test-failfast: ## Запустить тесты с остановкой на первой ошибке
	@echo "$(GREEN)Запуск тестов с остановкой на первой ошибке...$(NC)"
	@go test -failfast ./...

# Тесты по модулям
test-keymanager: ## Запустить тесты keymanager
	@echo "$(GREEN)Запуск тестов keymanager...$(NC)"
	@go test ./unit/keymanager/...

test-crypto: ## Запустить тесты crypto
	@echo "$(GREEN)Запуск тестов crypto...$(NC)"
	@go test ./unit/crypto/...

test-codec: ## Запустить тесты codec
	@echo "$(GREEN)Запуск тестов codec...$(NC)"
	@go test ./unit/codec/...

test-vault: ## Запустить тесты vault
	@echo "$(GREEN)Запуск тестов vault...$(NC)"
	@go test ./unit/vault/...

test-blockchain: ## Запустить тесты blockchain
	@echo "$(GREEN)Запуск тестов blockchain...$(NC)"
	@go test ./unit/blockchain/...

# Линтинг
lint: ## Запустить линтер
	@echo "$(GREEN)Запуск линтера...$(NC)"
	@golangci-lint run ./...

lint-fix: ## Исправить ошибки линтера
	@echo "$(GREEN)Исправление ошибок линтера...$(NC)"
	@golangci-lint run --fix ./...

# Очистка
clean: ## Очистить временные файлы
	@echo "$(GREEN)Очистка временных файлов...$(NC)"
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@go clean -testcache

# Установка зависимостей
install-deps: ## Установить зависимости для тестирования
	@echo "$(GREEN)Установка зависимостей...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/golang/mock/mockgen@latest

# Генерация моков
generate-mocks: ## Генерировать моки
	@echo "$(GREEN)Генерация моков...$(NC)"
	@mockgen -source=../internal/blockchain/client.go -destination=mocks/blockchain/mock_client.go -package=mocks
	@mockgen -source=../internal/crypto/crypto.go -destination=mocks/crypto/mock_crypto.go -package=mocks

# Проверка качества кода
quality-check: lint test-cover ## Проверить качество кода
	@echo "$(GREEN)Проверка качества кода завершена$(NC)"

# CI/CD команды
ci-test: test-race test-cover ## Команды для CI/CD
	@echo "$(GREEN)CI тесты завершены$(NC)"

# Разработка
dev-test: test-unit test-cover-html ## Команды для разработки
	@echo "$(GREEN)Тесты разработки завершены$(NC)"

# Полная проверка
full-check: clean install-deps generate-mocks quality-check ## Полная проверка проекта
	@echo "$(GREEN)Полная проверка завершена$(NC)"

# Справка по переменным окружения
env-help: ## Показать переменные окружения для тестов
	@echo "$(GREEN)Переменные окружения для тестов:$(NC)"
	@echo ""
	@echo "$(YELLOW)Для тестов блокчейна:$(NC)"
	@echo "  TEST_RPC_ENDPOINT=http://localhost:8545"
	@echo "  TEST_CONTRACT_ADDRESS=0x..."
	@echo "  TEST_PRIVATE_KEY=0x..."
	@echo ""
	@echo "$(YELLOW)Для тестов с файлами:$(NC)"
	@echo "  TEST_DATA_DIR=/tmp/encryptkeep_test"
	@echo ""
	@echo "$(YELLOW)Уровень логирования:$(NC)"
	@echo "  TEST_LOG_LEVEL=debug"
	@echo ""
	@echo "$(YELLOW)Пример запуска с переменными:$(NC)"
	@echo "  TEST_RPC_ENDPOINT=http://localhost:8545 make test-blockchain"
