FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY proto ./proto
COPY services/playback-service/go.mod services/playback-service/go.sum ./services/playback-service/

WORKDIR /app/services/playback-service
RUN go mod download

WORKDIR /app
COPY services/playback-service ./services/playback-service

WORKDIR /app/services/playback-service
RUN go build -o /app/bin/playback-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/playback-service .
COPY --from=builder /app/services/playback-service/migrations ./migrations

EXPOSE 50053

CMD ["./playback-service"]
