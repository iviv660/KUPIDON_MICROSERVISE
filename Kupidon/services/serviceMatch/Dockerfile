# Этап сборки (builder)
FROM golang:1.23.3-alpine AS builder

WORKDIR /serviceMatch

# Установка сертификатов и необходимых зависимостей
RUN apk add --no-cache ca-certificates openssl

# Копируем файлы для работы с модулями Go
COPY go.mod go.sum ./

# Загружаем зависимости Go
RUN go mod tidy && go mod download

# Копируем весь проект
COPY . .

# Сборка приложения
RUN go build -o /serviceMatch/main ./cmd/main.go

# Этап миграции
FROM golang:1.23.3-alpine AS migrate

WORKDIR /serviceMatch

# Установка необходимых зависимостей и драйвера для PostgreSQL
RUN apk add --no-cache bash && \
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Проверка, что migrate доступен
RUN which migrate || echo "migrate not found"

# Финальный образ
FROM alpine:latest

WORKDIR /serviceMatch

# Устанавливаем сертификаты для SSL
RUN apk --no-cache add ca-certificates bash postgresql-client

# Копируем исполнимые файлы из предыдущих слоев
COPY --from=builder /serviceMatch/main ./main
COPY --from=migrate /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /serviceMatch/migrations /migrations

# Устанавливаем переменные окружения
ENV DATABASE_URL="" \
    KAFKA_URL="" \
    KAFKA_LIKE_TOPIC="" \
    SERVICE_MATCH=""

# Открываем порт для приложения
EXPOSE 8081

# Запуск миграций и приложения
ENTRYPOINT ["sh", "-c", "migrate -path /migrations -database \"$DATABASE_URL\" -verbose up && ./main"]
