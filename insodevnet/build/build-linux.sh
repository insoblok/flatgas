#!/bin/bash
set -e

echo "🔧 Building Flatgas Linux binaries..."

# Determine repo root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
BIN_DIR="$SCRIPT_DIR/bin"

echo "📂 Script location: $SCRIPT_DIR"
echo "📂 Repo root: $REPO_ROOT"
echo "📂 Output bin dir: $BIN_DIR"

mkdir -p "$BIN_DIR"

TARGETS=(abigen bootnode clef ethkey evm faucet geth rlpdump)

for TARGET in "${TARGETS[@]}"; do
  SRC="$REPO_ROOT/cmd/$TARGET"
  OUT="$BIN_DIR/$TARGET"

  if [ -d "$SRC" ]; then
    echo "🔨 Building $TARGET..."
    GOOS=linux GOARCH=amd64 go build -o "$OUT" "$SRC"
  else
    echo "❌ Skipping $TARGET: source directory not found ($SRC)"
  fi
done

echo "🎉 Done. Binaries in $BIN_DIR"

# After geth is successfully built
if  [ -f "$OUT_BIN_DIR/geth" ]; then
  echo "🔁 Renaming geth to inso..."
  mv "$OUT_BIN_DIR/geth" "$OUT_BIN_DIR/inso"
fi

