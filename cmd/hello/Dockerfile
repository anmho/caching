FROM golang:1.22.2 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o hello ./cmd/hello

CMD ["./hello"]