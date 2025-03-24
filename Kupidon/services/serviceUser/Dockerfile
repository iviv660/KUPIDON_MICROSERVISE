# Сборка образа для приложения
FROM golang:1.23.3-alpine AS builder  

WORKDIR /serviceUser

# Устанавливаем сертификаты для работы с TLS
RUN apk add --no-cache ca-certificates openssl

# Копируем только go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . .

# Строим приложение
RUN go build -o /serviceUser/main ./cmd/main.go  # Исправлен путь сборки

# Сборка образа для миграций
FROM golang:1.23.3-alpine AS migrate

WORKDIR /serviceUser

# Устанавливаем bash и миграции
RUN apk add --no-cache bash && \
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Копируем миграции
COPY ./migrations ./migrations

# Финальный образ
FROM alpine:latest

WORKDIR /serviceUser

# Устанавливаем сертификаты
RUN apk --no-cache add ca-certificates bash

# Копируем приложение и миграции из предыдущих этапов
COPY --from=builder /serviceUser/main ./main
COPY --from=migrate /serviceUser/migrations /serviceUser/migrations
COPY --from=migrate /go/bin/migrate /usr/local/bin/migrate

# Устанавливаем переменные окружения
ENV USER_SERVICE="" \
    DATABASE_URL="" \
    MINIO_ENDPOINT="" \
    MINIO_ROOT_USER="" \
    MINIO_ROOT_PASSWORD="" \
    S3_BUCKET="" \
    REDIS_ADDR=""

EXPOSE 8080

# Запуск миграций и приложения
ENTRYPOINT ["sh", "-c", "migrate -path ./migrations -database \"$DATABASE_URL\" -verbose up && ./main"]
