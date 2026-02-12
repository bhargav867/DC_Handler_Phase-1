# Base stage
FROM golang:1.21-alpine AS base

WORKDIR /app

# Install basic utilities
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Development stage
FROM base AS development

COPY . .

ENV GOFLAGS="-v"
EXPOSE 8080

# Install air for hot-reload (optional, but useful)
RUN go install github.com/cosmtrek/air@latest

CMD ["go", "run", "main.go"]

# Build stage
FROM base AS builder

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dc-handler .

# Production stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/dc-handler .
COPY --from=builder /app/.env.example ./.env.example

EXPOSE 8080

CMD ["./dc-handler"]
