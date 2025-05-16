#!/bin/bash

set -e

REPO="ammar-ahmed22/chlog"
VERSION=${1:-"latest"}

# Detect OS
OS=$(uname -s)
case "$OS" in
    Linux*)     OS="linux" ;;
    Darwin*)    OS="darwin" ;;
    CYGWIN*|MINGW32*|MSYS*|MINGW*) OS="windows" ;;
    *)          echo "Unsupported OS: $OS" && exit 1 ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)     ARCH="x86_64" ;;
    arm64)      ARCH="arm64" ;;
    i386)       ARCH="i386" ;;
    *)          echo "Unsupported architecture: $ARCH" && exit 1 ;;
esac

echo "Detected: $OS $ARCH"

# Resolve version
if [ "$VERSION" == "latest" ]; then
  VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"v?([^"]+)".*/\1/')
fi

# Download and extract
FILENAME="chlog_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/v$VERSION/$FILENAME"

echo "Downloading $URL..."
curl -L "$URL" -o "$FILENAME"
tar -xzf "$FILENAME"

rm "$FILENAME"

# Move to /usr/local/bin
chmod +x chlog
sudo mv chlog /usr/local/bin/chlog

echo "✅ chlog v$VERSION installed to /usr/local/bin/chlog"
