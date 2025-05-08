#!/bin/bash
set -e

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
REPO_ROOT=$(cd "$SCRIPT_DIR/../.." && pwd)
IMAGE_NAME="insoblok/inso-node:dev-amd64"

echo "🐳 Building Docker image for linux/amd64..."
echo "📂 Script dir: $SCRIPT_DIR"
echo "📂 Repo root: $REPO_ROOT"
echo "🖼️ Image name: $IMAGE_NAME"

docker buildx build --platform linux/amd64 --load -t "$IMAGE_NAME" "$SCRIPT_DIR"

echo "✅ Image built for linux/amd64: $IMAGE_NAME"
