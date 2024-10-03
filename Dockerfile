# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

# Install git (required for some Go modules)
RUN apk update && apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
# Replace 'main.go' with the entry point of your application if different
RUN go build -o rto-attendance-tracker main.go

# Stage 2: Create a minimal image with the built binary
FROM alpine:latest

# Install necessary packages (e.g., certificates for HTTPS)
RUN apk update && apk add --no-cache ca-certificates && rm -rf /var/cache/apk/*

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
# Replace '8080' with your application's port if different
EXPOSE 8080

# Define environment variables if needed
# ENV DB_PATH=/app/data/db.sqlite3

# Set the entrypoint to run the binary
ENTRYPOINT ["./rto-attendance-tracker"]
