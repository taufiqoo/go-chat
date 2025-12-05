#!/bin/bash

echo "Building application..."

# Build for Linux (for Docker/GCP)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/chat-app-linux ./cmd/server

# Build for current OS
go build -o bin/chat-app ./cmd/server

echo "Build completed!"
echo "Linux binary: bin/chat-app-linux"
echo "Local binary: bin/chat-app"