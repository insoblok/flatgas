#!/bin/bash
set -e

# Resolve repository root from the script location
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
BUILD_DIR="$SCRIPT_DIR/bin"

echo "🔧 Building Flatgas Linux binaries..."
echo "📂 Script location: $SCRIPT_DIR"
echo "📂 Repo root: $REPO_ROOT"
echo "📂 Output bin dir: $BUILD_DIR"

# Clean bin directory
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# List of commands to build
TARGETS=("geth")

for TARGET in "${TARGETS[@]}"; do
  CMD_DIR="$REPO_ROOT/cmd/$TARGET"
  if [ ! -d "$CMD_DIR" ]; then
    echo "❌ Skipping $TARGET: source directory not found ($CMD_DIR)"
    continue
  fi

  echo "🛠️  Building $TARGET from $CMD_DIR..."
  GOOS=linux GOARCH=amd64 go build -o "$BUILD_DIR/$TARGET" "$CMD_DIR"
  echo "✅ Built $TARGET → $BUILD_DIR/$TARGET"
done

echo "🎉 Done. Binaries in $BUILD_DIR"
