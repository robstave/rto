# Stage 1: Build the Go application with cgo enabled
FROM golang:1.23-alpine AS builder

# Install necessary build tools for cgo
RUN apk add --no-cache git gcc musl-dev

# Set build arguments
ARG CACHEBUST=1
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o rto-attendance-tracker cmd/main/main.go

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Install necessary runtime packages
RUN apk add --no-cache \
    ca-certificates \
    sqlite-libs

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/rto-attendance-tracker .

# Copy static files and templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates

# Create a directory for data
RUN mkdir -p /app/data

# Expose the port the app runs on
EXPOSE 8080

# Define environment variables
ENV DB_PATH=/app/data/db.sqlite3

COPY --from=builder /app/static ./static

# Add a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Change ownership of the app directory
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

# Set the entrypoint to run the binary
ENTRYPOINT ["./rto-attendance-tracker"]
