#!/bin/bash
# Production build script for go-adapt

set -e  # Exit on error

echo "Building go-adapt for production..."

# Clean previous builds
rm -f go-adapt go-adapt-linux

# Build for Linux (your server)
echo "Building Linux binary..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o go-adapt-linux main.go

echo ""
echo "✓ Build complete!"
echo ""
echo "Production binary: go-adapt-linux"
echo "Size: $(du -h go-adapt-linux | cut -f1)"
echo ""
echo "To deploy:"
echo "  1. Upload go-adapt-linux and .env to your server"
echo "  2. On server: chmod +x go-adapt-linux"
echo "  3. Run with: GIN_MODE=release ./go-adapt-linux"
echo ""
echo "Optional: Build local binary for testing"
read -p "Build local binary? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Building local binary..."
    go build -o go-adapt main.go
    echo "✓ Local binary: go-adapt"
fi
