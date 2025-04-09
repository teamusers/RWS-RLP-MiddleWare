#!/bin/sh
set -e

# Configure parameters
PROJECT_DIR="/home/rlp-middleware-api"         # Directory where the code is located
OUTPUT_DIR="/app/rlp-middleware-api"           # Build output directory
BINARY_NAME="rlp-middleware-api"               # Name of the generated binary file

echo "Switch to the project directory: $PROJECT_DIR"
cd "$PROJECT_DIR"

echo "Pull the latest code..."
git pull

echo "Start compiling the Go program..."
export CGO_ENABLED=1
export GOOS=linux
export GOARCH=amd64

# Compile the project and output the binary file to the specified directory
go build -o "$OUTPUT_DIR/$BINARY_NAME" main.go
echo "Compilation successful, the binary file is located at $OUTPUT_DIR/$BINARY_NAME"