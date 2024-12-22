# Build stage
FROM golang:1.22-alpine AS builder

# Install git and build dependencies
RUN apk add --no-cache git build-base

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Final stage
FROM alpine:3.19

# Add non root user
RUN adduser -D -g '' appuser

# Install certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy any additional configuration files if needed
COPY --from=builder /app/internal/config/config.* ./internal/config/

# Use non root user
USER appuser

# Expose port (adjust if needed based on your config)
EXPOSE 8080

# Run the application
CMD ["./server"] 