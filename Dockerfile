# Build the Go application using the official Go image
FROM golang:1.22 AS builder
WORKDIR /app

# Copy go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . ./

# Build the Go application with static linking
RUN CGO_ENABLED=0 go build -o admin-auth-service ./config/main.go

# Start from a minimal base image
FROM debian:bullseye-slim

# Install CA certificates and other necessary dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Create config directory
RUN mkdir -p /app/config

# Copy the binary and .env file from the builder stage
COPY --from=builder /app/admin-auth-service /app/admin-auth-service
COPY --from=builder /app/config/.env /app/config/.env

# Expose the port
EXPOSE 50051

# Define the command to run the binary
CMD ["/app/admin-auth-service"]