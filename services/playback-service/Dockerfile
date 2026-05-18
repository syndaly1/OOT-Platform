FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o playback-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/playback-service .
COPY --from=builder /app/migrations ./migrations

EXPOSE 50053

CMD ["./playback-service"]