# Этап 1: Сборка приложения
FROM golang:1.24 AS builder

WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Генерируем Swagger-документацию
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init --dir ./cmd,./internal/api,./internal/models,./internal/repository,./internal/services,./internal/config,./internal/logger --parseDependency

# Компилируем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o persons-api ./cmd

# Этап 2: Создание минимального образа для запуска
FROM alpine:latest

WORKDIR /app

# Копируем скомпилированный бинарный файл из этапа сборки
COPY --from=builder /app/persons-api .

# Открываем порт, на котором будет работать приложение
EXPOSE 8080

# Запускаем приложение
CMD ["./persons-api"]