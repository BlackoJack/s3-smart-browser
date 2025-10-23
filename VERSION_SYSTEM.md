# Система автоматических версий

## Обзор

Теперь версия приложения автоматически извлекается из git тегов и коммитов, что исключает необходимость ручного обновления версии в коде.

## Как это работает

### 1. Автоматическое определение версии
- **Релизные версии**: берутся из последнего git тега (например, `1.0.8`)
- **Dev версии**: формируются как `тег-количество-коммитов-хеш` (например, `1.0.8-1-g2880d5c-dirty`)
- **Коммит**: полный хеш текущего коммита
- **Время сборки**: автоматически устанавливается при сборке

### 2. Переменные сборки
При сборке через `ldflags` устанавливаются следующие переменные:
- `s3-smart-browser/internal/version.Version` - версия из git
- `s3-smart-browser/internal/version.GitCommit` - хеш коммита
- `s3-smart-browser/internal/version.BuildTime` - время сборки

## Способы сборки

### 1. Через Makefile (рекомендуется)
```bash
# Сборка с автоматической версией
make build

# Запуск в режиме разработки
make run

# Полная сборка для релиза
make release

# Показать информацию о версии
make version
```

### 2. Через скрипт build.sh
```bash
# Сборка
./build.sh

# Сборка и запуск
./build.sh --run
```

### 3. Ручная сборка
```bash
VERSION=$(git describe --tags --always --dirty)
GIT_COMMIT=$(git rev-parse HEAD)
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

go build \
    -ldflags "-X s3-smart-browser/internal/version.Version=$VERSION -X s3-smart-browser/internal/version.GitCommit=$GIT_COMMIT -X s3-smart-browser/internal/version.BuildTime=$BUILD_TIME" \
    -o build/s3-smart-browser \
    cmd/server/main.go
```

## API версии

### Endpoint: `/api/version`

Возвращает JSON с информацией о версии:
```json
{
  "version": "1.0.8-1-g2880d5c-dirty",
  "git_commit": "2880d5c7571afa9de1beec990e3168b81f4f9403",
  "build_time": "2025-10-23_16:49:20",
  "go_version": "go1.25.1"
}
```

## Отображение в UI

В футере приложения отображается:
- **Название приложения**: "S3 Smart Browser"
- **Версия**: автоматически загружается из API (например, "v1.0.8-1-g2880d5c-dirty")
- **Статистика**: количество файлов и статус подключения
- **Ссылки**: GitHub и Telegram

## Рабочий процесс

### Для разработки
1. Работайте в обычном режиме
2. Используйте `make run` для запуска с dev версией
3. Версия будет автоматически включать информацию о коммите

### Для релиза
1. Создайте git тег: `git tag v1.0.9`
2. Запустите `make release`
3. Версия будет автоматически установлена в `1.0.9`

### Для CI/CD

#### GitHub Actions
Автоматически работает с обновленными workflow:
- **Релизы**: версия берется из тега (`github.ref_name`)
- **Push/PR**: версия формируется из git describe
- **Docker**: переменные версии передаются через build-args

#### Docker сборка
```bash
# Локальная сборка с версией
make docker-build

# Ручная сборка
docker build \
  --build-arg VERSION=$(git describe --tags --always --dirty) \
  --build-arg GIT_COMMIT=$(git rev-parse HEAD) \
  --build-arg BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S') \
  -t s3-smart-browser:latest .
```

#### Переменные окружения
```bash
VERSION=$(git describe --tags --always --dirty)
GIT_COMMIT=$(git rev-parse HEAD)
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
# Используйте в вашем пайплайне
```

## Преимущества

✅ **Автоматизация**: не нужно вручную обновлять версию в коде  
✅ **Трассируемость**: каждый билд содержит информацию о коммите  
✅ **Гибкость**: поддержка как релизных, так и dev версий  
✅ **Интеграция**: легко интегрируется с CI/CD системами  
✅ **Отладка**: полная информация о сборке для диагностики  

## Примеры версий

- `1.0.8` - релизная версия
- `1.0.8-1-g2880d5c` - версия после релиза (1 коммит после тега)
- `1.0.8-1-g2880d5c-dirty` - версия с незакоммиченными изменениями
- `dev-abc12345` - версия без тегов (dev + короткий хеш)
