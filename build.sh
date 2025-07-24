#!/usr/bin/env bash

# Build script for BBS Sabacc
# This script builds the game for various platforms

# Default to current platform if no argument provided
TARGET_OS=${1:-$(uname -s | tr '[:upper:]' '[:lower:]')}
TARGET_ARCH=${2:-amd64}

if [ "$TARGET_OS" = "darwin" ]; then
    TARGET_OS="darwin"
elif [ "$TARGET_OS" = "linux" ]; then
    TARGET_OS="linux"
elif [ "$TARGET_OS" = "windows" ]; then
    TARGET_OS="windows"
fi

echo "Building BBS Sabacc for $TARGET_OS/$TARGET_ARCH..."

# Clean any previous builds
rm -f sabacc sabacc.exe

# Set build environment
export CGO_ENABLED=0
export GOOS=$TARGET_OS
export GOARCH=$TARGET_ARCH

# Download dependencies
echo "Downloading dependencies..."
go mod tidy

# Build the card database first
echo "Building card database..."
go run cmd/build-cards/main.go

# Set executable name based on OS
EXECUTABLE="sabacc"
if [ "$TARGET_OS" = "windows" ]; then
    EXECUTABLE="sabacc.exe"
fi

# Build the executable
echo "Compiling..."
go build -ldflags="-s -w" -o $EXECUTABLE .

if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo "Target: $TARGET_OS/$TARGET_ARCH"
    echo "Executable: $EXECUTABLE"
    echo "Size: $(du -h $EXECUTABLE | cut -f1)"
    echo ""
    
    if [ "$TARGET_OS" = "darwin" ]; then
        echo "macOS Installation:"
        echo "1. Copy '$EXECUTABLE' to desired location"
        echo "2. Run with: ./$EXECUTABLE -path [dropfile_path]"
        echo "3. For testing: echo -e '2\\n8\\n38400\\nTest BBS\\n1\\nTest Player\\nTestUser\\n100\\n90\\n0\\n1' > door32.sys"
        echo "4. Then run: ./$EXECUTABLE -path ./"
    elif [ "$TARGET_OS" = "linux" ]; then
        echo "Linux BBS Installation:"
        echo "1. Copy '$EXECUTABLE' to your BBS doors directory"
        echo "2. Configure your BBS to launch with: ./$EXECUTABLE -path [dropfile_path]"
        echo "3. Ensure the executable has proper permissions: chmod +x $EXECUTABLE"
    elif [ "$TARGET_OS" = "windows" ]; then
        echo "Windows Installation:"
        echo "1. Copy '$EXECUTABLE' to your BBS doors directory"
        echo "2. Configure your BBS to launch with: $EXECUTABLE -path [dropfile_path]"
    fi
    
    echo ""
    echo "Card database: sabacc_cards.bin ($(du -h sabacc_cards.bin 2>/dev/null | cut -f1 || echo 'not found'))"
    
else
    echo "Build failed!"
    exit 1
fi

echo ""
echo "Usage examples:"
echo "  Build for current platform: ./build.sh"
echo "  Build for Linux BBS:        ./build.sh linux"
echo "  Build for macOS:            ./build.sh darwin"
echo "  Build for Windows:          ./build.sh windows"
