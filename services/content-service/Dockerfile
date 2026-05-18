FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o content-service ./cmd

EXPOSE 50052

CMD ["./content-service"]