# Stage 1: Build the Go binary with CGO disabled
FROM golang:1.22.3-alpine AS builder

# Set the working directory in the builder stage
WORKDIR /app/lbe-api

# Copy go.mod and go.sum files first to leverage caching and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code into the container (including the config folder)
COPY . .

# Create a directory for binary output
RUN mkdir -p /app/bin

# Build the Go binary with CGO disabled, outputting it to /app/bin/lbe-api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/lbe-api main.go

# Stage 2: Create a minimal runtime container
FROM alpine:latest

# Set the working directory to /app/lbe-api (this ensures it is a directory)
WORKDIR /app/lbe-api

# Copy the built binary from the builder stage into this directory
COPY --from=builder /app/bin/lbe-api .

# If needed, copy the config folder from the builder stage (adjust if you want it inside this directory)
COPY --from=builder /app/lbe-api/config ./config

# Expose the port your application listens on
EXPOSE 18080

# Command to run the binary
CMD ["./lbe-api"]
