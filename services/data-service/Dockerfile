FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копирование go mod файлов
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN go build -o data-service ./cmd

# Финальный образ
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR app/

# Копирование бинарника
COPY --from=builder /app/data-service .

# Копирование миграций
COPY migrations ./migrations

# Экспорт порта
EXPOSE 8082

# Запуск приложения
CMD ["./data-service"]
