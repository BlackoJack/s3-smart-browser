#!/bin/bash

# Скрипт для автоматической сборки S3 Smart Browser с версией из git

set -e

# Цвета
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Функция для вывода сообщений
log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Проверка наличия git
if ! command -v git &> /dev/null; then
    error "Git не найден. Установите git для работы с версиями."
    exit 1
fi

# Получение информации о версии
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION=$(go version | cut -d' ' -f3)

log "Информация о сборке:"
echo "  Версия: $VERSION"
echo "  Коммит: $GIT_COMMIT"
echo "  Время сборки: $BUILD_TIME"
echo "  Go версия: $GO_VERSION"
echo ""

# Создание директории для сборки
BUILD_DIR="build"
mkdir -p "$BUILD_DIR"

# Сборка с переменными версии
log "Сборка приложения..."
go build \
    -ldflags "-X s3-smart-browser/internal/version.Version=$VERSION -X s3-smart-browser/internal/version.GitCommit=$GIT_COMMIT -X s3-smart-browser/internal/version.BuildTime=$BUILD_TIME" \
    -o "$BUILD_DIR/s3-smart-browser" \
    cmd/server/main.go

if [ $? -eq 0 ]; then
    log "Сборка успешно завершена!"
    log "Исполняемый файл: $BUILD_DIR/s3-smart-browser"
    
    # Показать размер файла
    FILE_SIZE=$(du -h "$BUILD_DIR/s3-smart-browser" | cut -f1)
    echo "  Размер: $FILE_SIZE"
    
    # Опционально запустить приложение
    if [ "$1" = "--run" ]; then
        log "Запуск приложения..."
        ./"$BUILD_DIR/s3-smart-browser"
    fi
else
    error "Ошибка при сборке!"
    exit 1
fi
