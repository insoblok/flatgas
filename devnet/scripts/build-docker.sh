#!/bin/bash
set -euo pipefail

TAG=${1:-flatgas-node}
DOCKERFILE=devnet/Dockerfile

echo "📦 Building Docker image: $TAG"
docker build -f "$DOCKERFILE" -t "$TAG" .
