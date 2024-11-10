FROM golang:1.22.2 AS builder
WORKDIR /app

ENV GOARCH=amd64

# Dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Build application
COPY . .
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH go build -o hello ./cmd/hello

FROM alpine:3.20.3

WORKDIR /app

COPY --from=builder /app/hello .

EXPOSE 8080
CMD ["./hello"]
