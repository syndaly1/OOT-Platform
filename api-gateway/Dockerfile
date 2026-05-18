FROM golang:1.25-alpine AS builder

WORKDIR /app

ENV GOWORK=off

COPY go.mod go.sum ./
COPY proto ./proto
COPY api-gateway ./api-gateway

RUN go mod download
RUN go build -o /app/api-gateway ./api-gateway

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api-gateway .

EXPOSE 8080

CMD ["./api-gateway"]