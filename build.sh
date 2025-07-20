#!/bin/bash

# Build script for BBS Sabacc
# This script builds the game for Linux BBS systems

echo "Building BBS Sabacc..."

# Clean any previous builds
rm -f sabacc

# Set build environment
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

# Download dependencies
echo "Downloading dependencies..."
go mod tidy

# Build the executable
echo "Compiling..."
go build -ldflags="-s -w" -o sabacc .

if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo "Executable: sabacc"
    echo "Size: $(du -h sabacc | cut -f1)"
    echo ""
    echo "Installation:"
    echo "1. Copy 'sabacc' to your BBS doors directory"
    echo "2. Configure your BBS to launch with: ./sabacc -path [dropfile_path]"
    echo "3. Ensure the executable has proper permissions: chmod +x sabacc"
else
    echo "Build failed!"
    exit 1
fi