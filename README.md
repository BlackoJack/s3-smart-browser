# S3 Smart Browser

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Helm-326CE5?style=for-the-badge&logo=kubernetes)](https://kubernetes.io/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![GitHub Release](https://img.shields.io/github/v/release/BlackoJack/s3-smart-browser?style=for-the-badge)](https://github.com/BlackoJack/s3-smart-browser/releases)

> 🚀 **Современный веб-браузер для Amazon S3 с интуитивным интерфейсом**

S3 Smart Browser — это легковесное веб-приложение на Go, которое предоставляет удобный веб-интерфейс для просмотра и скачивания файлов из Amazon S3 хранилища. Приложение поддерживает навигацию по папкам, скачивание файлов и работает с любыми S3-совместимыми хранилищами.

## ✨ Основные возможности

- 📁 **Навигация по папкам** - интуитивный интерфейс для просмотра структуры файлов
- 📥 **Скачивание файлов** - прямое скачивание через presigned URLs
- 🔍 **Поиск и фильтрация** - быстрый поиск по файловой системе
- 📱 **Адаптивный дизайн** - работает на всех устройствах
- 🚀 **Высокая производительность** - асинхронная обработка запросов
- 🔒 **Безопасность** - поддержка IAM ролей и временных токенов
- 🐳 **Контейнеризация** - готовые Docker образы
- ☸️ **Kubernetes** - Helm чарты для развертывания

## 🏗️ Архитектура

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Browser   │────│  Go Web Server  │────│   S3 Storage    │
│                 │    │                 │    │                 │
│ • HTML/CSS/JS   │    │ • HTTP Handlers │    │ • File Storage  │
│ • File Browser  │    │ • S3 Client     │    │ • Presigned URLs│
│ • Download UI   │    │ • Config Mgmt   │    │ • Directory API │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 Быстрый старт

### Предварительные требования

- Docker и Docker Compose
- Доступ к S3 хранилищу (AWS S3, MinIO, или другое S3-совместимое)
- AWS Access Key ID и Secret Access Key

### Запуск с Docker Compose

1. **Клонируйте репозиторий:**
```bash
git clone https://github.com/BlackoJack/s3-smart-browser.git
cd s3-smart-browser
```

2. **Создайте файл окружения:**
```bash
cp env.example .env
```

3. **Настройте переменные окружения в `.env`:**
```env
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_BUCKET=your-bucket-name
AWS_REGION=us-east-1
BASE_DIRECTORY=optional/base/path
```

4. **Запустите приложение:**
```bash
docker-compose up -d
```

5. **Откройте браузер:**
```
http://localhost:8080
```

### Запуск с Docker

```bash
docker run -d \
  --name s3-smart-browser \
  -p 8080:8080 \
  -e AWS_ACCESS_KEY_ID=your_access_key \
  -e AWS_SECRET_ACCESS_KEY=your_secret_key \
  -e AWS_BUCKET=your-bucket-name \
  -e AWS_REGION=us-east-1 \
  ghcr.io/blackojack/s3-smart-browser:latest
```

## ☸️ Развертывание в Kubernetes

### Использование Helm

1. **Добавьте репозиторий:**
```bash
helm repo add s3-smart-browser https://blackojack.github.io/s3-smart-browser
helm repo update
```

2. **Установите чарт:**
```bash
helm install s3-smart-browser s3-smart-browser/s3-smart-browser \
  --set config.aws.bucket=your-bucket-name \
  --set config.aws.region=us-east-1 \
  --set secrets.awsAccessKeyId=your_access_key \
  --set secrets.awsSecretAccessKey=your_secret_key
```

3. **Или используйте values.yaml:**
```bash
helm install s3-smart-browser s3-smart-browser/s3-smart-browser -f helm-values-example.yaml
```

### Пример values.yaml для Helm

```yaml
replicaCount: 2

image:
  repository: ghcr.io/blackojack/s3-smart-browser
  tag: latest
  pullPolicy: IfNotPresent

config:
  aws:
    region: "us-east-1"
    bucket: "my-s3-bucket"
    endpoint: ""  # Для MinIO: "http://minio:9000"
  server:
    port: 8080
  app:
    baseDirectory: "documents"  # Опциональная базовая папка

secrets:
  awsAccessKeyId: "AKIAIOSFODNN7EXAMPLE"
  awsSecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

ingress:
  enabled: true
  className: "nginx"
  hosts:
    - host: s3-browser.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: s3-browser-tls
      hosts:
        - s3-browser.example.com

resources:
  limits:
    cpu: 500m
    memory: 128Mi
  requests:
    cpu: 100m
    memory: 64Mi
```

## ⚙️ Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию | Обязательная |
|------------|----------|--------------|--------------|
| `AWS_ACCESS_KEY_ID` | AWS Access Key ID | - | ✅ |
| `AWS_SECRET_ACCESS_KEY` | AWS Secret Access Key | - | ✅ |
| `AWS_BUCKET` | Имя S3 bucket | - | ✅ |
| `AWS_REGION` | AWS регион | `us-east-1` | ❌ |
| `AWS_ENDPOINT` | Кастомный S3 endpoint | - | ❌ |
| `SERVER_PORT` | Порт сервера | `8080` | ❌ |
| `BASE_DIRECTORY` | Базовая папка в bucket | - | ❌ |

### Поддержка S3-совместимых хранилищ

Приложение работает с любыми S3-совместимыми хранилищами:

- **Amazon S3** - стандартная конфигурация
- **MinIO** - используйте `AWS_ENDPOINT=http://minio:9000`
- **DigitalOcean Spaces** - используйте соответствующий endpoint
- **Wasabi** - используйте Wasabi endpoint
- **Другие S3-совместимые** - настройте endpoint

## 🔧 Разработка

### Локальная разработка

1. **Установите Go 1.25+:**
```bash
go version
```

2. **Клонируйте репозиторий:**
```bash
git clone https://github.com/BlackoJack/s3-smart-browser.git
cd s3-smart-browser
```

3. **Установите зависимости:**
```bash
go mod download
```

4. **Настройте переменные окружения:**
```bash
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_BUCKET=your-bucket-name
export AWS_REGION=us-east-1
```

5. **Запустите приложение:**
```bash
go run cmd/server/main.go
```

### Сборка Docker образа

```bash
docker build -t s3-smart-browser .
```

### Структура проекта

```
s3-smart-browser/
├── cmd/server/           # Точка входа приложения
├── internal/
│   ├── config/          # Конфигурация
│   ├── handlers/        # HTTP обработчики
│   ├── s3/             # S3 клиент
│   └── types/          # Типы данных
├── web/
│   ├── static/         # Статические файлы (CSS, JS)
│   └── templates/      # HTML шаблоны
├── charts/             # Helm чарты
├── docker-compose.yaml # Docker Compose конфигурация
└── Dockerfile          # Docker образ
```

## 📊 API Endpoints

| Endpoint | Метод | Описание |
|----------|-------|----------|
| `/` | GET | Главная страница |
| `/api/list` | GET | Список файлов и папок |
| `/api/download` | GET | Скачивание файла |

### Примеры API запросов

**Получить список файлов:**
```bash
curl "http://localhost:8080/api/list?path=/documents"
```

**Скачать файл:**
```bash
curl -O "http://localhost:8080/api/download?file=/documents/file.pdf"
```

📖 **Подробные примеры API**: см. [API_EXAMPLES.md](API_EXAMPLES.md)

📦 **Примеры развертывания**: см. [DEPLOYMENT_EXAMPLES.md](DEPLOYMENT_EXAMPLES.md)

❓ **Часто задаваемые вопросы**: см. [FAQ.md](FAQ.md)

📋 **История изменений**: см. [CHANGELOG.md](CHANGELOG.md)

📢 **Информация о релизах**: см. [RELEASE_NOTES.md](RELEASE_NOTES.md)

☸️ **Helm Chart документация**: см. [HELM_CHART.md](HELM_CHART.md)

## 🛡️ Безопасность

- ✅ **Presigned URLs** - безопасное скачивание без передачи данных через сервер
- ✅ **IAM роли** - поддержка временных токенов и ролей
- ✅ **HTTPS** - поддержка SSL/TLS
- ✅ **Non-root контейнер** - запуск от непривилегированного пользователя
- ✅ **Security contexts** - настройки безопасности для Kubernetes

🔒 **Политика безопасности**: см. [SECURITY.md](SECURITY.md)

## 📈 Мониторинг и логирование

### Health Check

Приложение предоставляет health check endpoint:
```bash
curl http://localhost:8080/
```

### Логирование

Приложение использует стандартное логирование Go с уровнями:
- INFO - общая информация
- ERROR - ошибки
- FATAL - критические ошибки

## 🤝 Вклад в проект

Мы приветствуем вклад в развитие проекта! 

📖 **Подробные инструкции**: см. [CONTRIBUTING.md](CONTRIBUTING.md)

🛠️ **Руководство разработчика**: см. [DEVELOPMENT.md](DEVELOPMENT.md)

### Быстрый старт:
1. Форкните репозиторий
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Зафиксируйте изменения (`git commit -m 'Add amazing feature'`)
4. Отправьте в branch (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📝 Лицензия

Этот проект лицензирован под MIT License - см. файл [LICENSE](LICENSE) для деталей.

## 👨‍💻 Автор

**Raenko Ivan** - [@devopsnsk](https://t.me/devopsnsk)

- GitHub: [@BlackoJack](https://github.com/BlackoJack)
- Email: devopsnskru@gmail.com

## 🙏 Благодарности

- [AWS SDK for Go](https://github.com/aws/aws-sdk-go-v2) - за отличную поддержку S3
- [Helm](https://helm.sh/) - за упрощение развертывания в Kubernetes
- [Docker](https://www.docker.com/) - за контейнеризацию

## 📞 Поддержка

Если у вас есть вопросы или проблемы:

- 🐛 **Баг-репорты**: [GitHub Issues](https://github.com/BlackoJack/s3-smart-browser/issues)
- 💬 **Обсуждения**: [GitHub Discussions](https://github.com/BlackoJack/s3-smart-browser/discussions)
- 📧 **Email**: devopsnskru@gmail.com
- 💬 **Telegram**: [@devopsnsk](https://t.me/devopsnsk)

---

⭐ **Если проект был полезен, поставьте звезду!** ⭐
