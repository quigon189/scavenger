FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копирование go mod файлов
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN go build -o scavenger ./cmd

# Финальный образ
FROM alpine:latest
WORKDIR app/

# Копирование бинарника
COPY --from=builder /app/scavenger .

# Экспорт порта
EXPOSE ${PORT:-8080}

# Запуск приложения
CMD ["./scavenger"]
