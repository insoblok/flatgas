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

TARGETS=(abidump  abigen   blsync   clef     devp2p   era      ethkey   evm      geth     rlpdump    workload)

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

# After geth is successfully built
if  [ -f "$BIN_DIR/geth" ]; then
  echo "🔁 Renaming geth to inso..."
  mv "$BIN_DIR/geth" "$BIN_DIR/inso"
fi

echo "🎉 Done. Binaries in $BIN_DIR"
