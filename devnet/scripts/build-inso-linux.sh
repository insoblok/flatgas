#!/bin/bash
set -e

echo "🛠️  Building Linux AMD64 binary for Docker..."

mkdir -p devnet/build

GOOS=linux GOARCH=amd64 go build -o devnet/build/inso ./cmd/geth

echo "✅ Build complete: devnet/build/inso"
file devnet/build/inso
