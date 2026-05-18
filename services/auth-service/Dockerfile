FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY proto ./proto
COPY services/auth-service ./services/auth-service

WORKDIR /app/services/auth-service

RUN go mod download
RUN go build -o /app/auth-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth-service .
COPY --from=builder /app/services/auth-service/migrations ./migrations

EXPOSE 50051

CMD ["./auth-service"]