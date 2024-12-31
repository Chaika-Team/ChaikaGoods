# Use an official Golang image with version 1.23.4 as a base image
FROM golang:1.23.4-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create an app directory
WORKDIR /app

# Install dependencies for the build
RUN apk add --no-cache git

# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . ./

# Build the application
RUN go build -o main ./cmd/main.go

# Use a minimal image for running the application
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy configuration file через volume (не копируем сюда)

# Expose the application port
EXPOSE ${APP_PORT:-8080}

# Run the binary
CMD ["./main", "-config=/config/config.yml"]
