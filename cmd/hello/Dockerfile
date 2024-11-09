FROM golang:1.22.2 AS builder
WORKDIR /app

# Dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Build application
COPY . .
RUN go build -o hello ./cmd/hello

FROM alpine:3.20.3
WORKDIR /app
COPY --from=builder /app/hello .
CMD ["./hello"]

