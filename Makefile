# S3 Smart Browser Makefile

# Переменные
APP_NAME := s3-smart-browser
BUILD_DIR := build
VERSION := $(shell git describe --tags --always --dirty)
GIT_COMMIT := $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X s3-smart-browser/internal/version.Version=$(VERSION) -X s3-smart-browser/internal/version.GitCommit=$(GIT_COMMIT) -X s3-smart-browser/internal/version.BuildTime=$(BUILD_TIME)"

# Цвета для вывода
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
NC := \033[0m # No Color

.PHONY: help build run clean test docker-build docker-run dev install

# Помощь
help: ## Показать справку
	@echo "$(GREEN)S3 Smart Browser - Доступные команды:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(GREEN)Текущая версия:$(NC) $(VERSION)"
	@echo "$(GREEN)Git commit:$(NC) $(GIT_COMMIT)"

# Сборка
build: ## Собрать приложение
	@echo "$(GREEN)Сборка $(APP_NAME)...$(NC)"
	@echo "$(YELLOW)Версия:$(NC) $(VERSION)"
	@echo "$(YELLOW)Коммит:$(NC) $(GIT_COMMIT)"
	@echo "$(YELLOW)Время сборки:$(NC) $(BUILD_TIME)"
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) cmd/server/main.go
	@echo "$(GREEN)Сборка завершена: $(BUILD_DIR)/$(APP_NAME)$(NC)"

# Запуск в режиме разработки
run: ## Запустить приложение в режиме разработки
	@echo "$(GREEN)Запуск $(APP_NAME) в режиме разработки...$(NC)"
	@echo "$(YELLOW)Версия:$(NC) $(VERSION)"
	go run $(LDFLAGS) cmd/server/main.go

# Запуск собранного приложения
run-build: build ## Собрать и запустить приложение
	@echo "$(GREEN)Запуск собранного приложения...$(NC)"
	./$(BUILD_DIR)/$(APP_NAME)

# Очистка
clean: ## Очистить файлы сборки
	@echo "$(YELLOW)Очистка...$(NC)"
	rm -rf $(BUILD_DIR)
	go clean

# Тесты
test: ## Запустить тесты
	@echo "$(GREEN)Запуск тестов...$(NC)"
	go test -v ./...

# Docker сборка
docker-build: ## Собрать Docker образ
	@echo "$(GREEN)Сборка Docker образа...$(NC)"
	@echo "$(YELLOW)Версия:$(NC) $(VERSION)"
	@echo "$(YELLOW)Коммит:$(NC) $(GIT_COMMIT)"
	@echo "$(YELLOW)Время сборки:$(NC) $(BUILD_TIME)"
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		-t $(APP_NAME):$(VERSION) .
	docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest

# Docker запуск
docker-run: docker-build ## Собрать и запустить Docker контейнер
	@echo "$(GREEN)Запуск Docker контейнера...$(NC)"
	docker run -p 8080:8080 --env-file env.example $(APP_NAME):$(VERSION)

# Тестирование Docker версии
docker-test-version: docker-build ## Тестировать версию в Docker контейнере
	@echo "$(GREEN)Тестирование версии в Docker...$(NC)"
	docker run -d --name $(APP_NAME)-test -p 8080:8080 $(APP_NAME):$(VERSION)
	@sleep 3
	@echo "$(YELLOW)Тестирование API версии...$(NC)"
	@curl -s http://localhost:8080/api/version | jq . || echo "$(RED)Ошибка при тестировании API$(NC)"
	@docker stop $(APP_NAME)-test
	@docker rm $(APP_NAME)-test

# Установка в систему
install: build ## Установить приложение в систему
	@echo "$(GREEN)Установка $(APP_NAME) в систему...$(NC)"
	sudo cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/
	@echo "$(GREEN)Установка завершена!$(NC)"

# Показать информацию о версии
version: ## Показать информацию о версии
	@echo "$(GREEN)Информация о версии:$(NC)"
	@echo "  Версия: $(VERSION)"
	@echo "  Коммит: $(GIT_COMMIT)"
	@echo "  Время сборки: $(BUILD_TIME)"
	@echo "  Go версия: $(shell go version)"

# Проверка зависимостей
deps: ## Проверить и обновить зависимости
	@echo "$(GREEN)Проверка зависимостей...$(NC)"
	go mod tidy
	go mod verify

# Форматирование кода
fmt: ## Форматировать код
	@echo "$(GREEN)Форматирование кода...$(NC)"
	go fmt ./...

# Линтинг
lint: ## Запустить линтер
	@echo "$(GREEN)Запуск линтера...$(NC)"
	golangci-lint run

# Полная сборка для релиза
release: clean deps fmt test build ## Полная сборка для релиза
	@echo "$(GREEN)Релизная сборка завершена!$(NC)"
	@echo "$(YELLOW)Файл:$(NC) $(BUILD_DIR)/$(APP_NAME)"
	@echo "$(YELLOW)Версия:$(NC) $(VERSION)"

# По умолчанию
.DEFAULT_GOAL := help
