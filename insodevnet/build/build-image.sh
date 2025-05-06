#!/bin/bash
set -e

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
REPO_ROOT=$(cd "$SCRIPT_DIR/../.." && pwd)
IMAGE_NAME="insoblok/inso-node:dev"

echo "🐳 Building Docker image..."
echo "📂 Script dir: $SCRIPT_DIR"
echo "📂 Repo root: $REPO_ROOT"
echo "🖼️ Image name: $IMAGE_NAME"

docker build -t "$IMAGE_NAME" "$SCRIPT_DIR"

echo "✅ Image built: $IMAGE_NAME"
