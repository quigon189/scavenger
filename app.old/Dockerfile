FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копирование go mod файлов
COPY go.work go.work.sum ./

COPY core/go.mod core/go.sum ./core/
COPY web/go.mod web/go.sum ./web/

RUN go work sync

COPY . .

# Сборка приложения
RUN go build -o scavenger ./web/cmd/main.go

# Финальный образ
FROM alpine:latest
WORKDIR app/

# Копирование бинарника
COPY --from=builder /app/scavenger .

# Экспорт порта
EXPOSE ${PORT:-8080}

# Запуск приложения
CMD ["./scavenger"]
