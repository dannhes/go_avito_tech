# Stage 1: Build the Go binary
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем бинарник
RUN go build -o pr_service ./cmd/service/main.go

# Stage 2: Минимальный runtime
FROM alpine:3.18

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/pr_service .

# Экспонируем порт
EXPOSE 8080

# Команда запуска
CMD ["./pr_service"]
