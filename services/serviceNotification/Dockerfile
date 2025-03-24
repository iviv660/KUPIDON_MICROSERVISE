FROM golang:1.23.3-alpine AS builder

WORKDIR /serviceNotification

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o /serviceNotification/cmd/main ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /serviceNotification

COPY --from=builder /serviceNotification/cmd/main /serviceNotification/cmd/main

ENV KAFKA_BROKER="" \
    KAFKA_LIKE_TOPIC=""\
    TELEGRAM_BOT_TOKEN=""\
    USER_SERVICE=""

CMD ["./cmd/main"]