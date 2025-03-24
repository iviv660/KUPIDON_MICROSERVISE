FROM golang:1.23.3-alpine AS builder

WORKDIR /serviceBot

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o /serviceBot/cmd/main ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /serviceBot

COPY --from=builder /serviceBot/cmd/main /serviceBot/cmd/main

# Убираем пробелы в переменных окружения
ENV TELEGRAM_BOT_TOKEN="" \
    USER_SERVICE="" \
    MATCH_SERVICE=""

CMD ["/serviceBot/cmd/main"]
