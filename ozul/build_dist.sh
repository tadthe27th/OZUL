#!/bin/sh
set -e
mkdir -p dist

echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o dist/ozul-windows-amd64.exe .

echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o dist/ozul-linux-amd64 .

echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o dist/ozul-darwin-amd64 .

echo "Build complete! Binaries are in the dist/ directory." 