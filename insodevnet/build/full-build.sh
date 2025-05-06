#!/bin/sh
set -e

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
BIN_DIR="$SCRIPT_DIR/bin"
echo "BinDir $BIN_DIR"
echo "🧹 Cleaning old binaries..."
rm -rf "$BIN_DIR"
mkdir -p "$BIN_DIR"

echo "🔨 Building Linux binaries..."
"$SCRIPT_DIR/build-linux.sh"

echo "🐳 Building Docker image..."
"$SCRIPT_DIR/build-image.sh"

echo "🚀 Verifying image by running: inso version"
docker run --rm insoblok/inso-node:dev  --version
