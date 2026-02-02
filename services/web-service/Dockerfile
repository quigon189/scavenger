FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копирование go mod файлов
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN go build -o web-service ./cmd

# Финальный образ
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR app/

# Копирование бинарника
COPY --from=builder /app/web-service .

# Экспорт порта
EXPOSE ${PORT:-8083}

# Запуск приложения
CMD ["./web-service"]
