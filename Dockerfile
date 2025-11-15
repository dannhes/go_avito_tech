FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o pr_service ./cmd/service/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/pr_service .

EXPOSE 8080

CMD ["./pr_service"]
