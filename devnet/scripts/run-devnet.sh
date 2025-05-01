#!/bin/bash
set -euo pipefail

# Defaults
DEFAULT_DATA_DIR="$(pwd)/devnet"
DEFAULT_INSOD="$(pwd)/build/bin/inso"

# Args
DATA_DIR="${1:-$DEFAULT_DATA_DIR}"
INSO_BIN="${2:-$DEFAULT_INSOD}"

# Check that inso binary exists
if [[ ! -x "$INSO_BIN" ]]; then
  echo "❌ Error: Cannot find or execute 'inso' at $INSO_BIN"
  exit 1
fi

echo "🚀 Starting inso node"
echo "📁 Data dir: $DATA_DIR"
echo "⚙️ Binary:   $INSO_BIN"

CMD=(
  "$INSO_BIN"
  --datadir "$DATA_DIR"
  --http
  --http.api admin,debug,eth,miner,net,txpool,web3
  --http.addr 0.0.0.0
  --http.port 8545
  --http.corsdomain "*"
  --ws
  --ws.api eth,net,web3,debug
  --ws.addr 0.0.0.0
  --ws.port 8546
  --ws.origins "*"
  --mine
  --nodiscover
  --networkid 12345
  --allow-insecure-unlock
)

echo "🧪 Command to run:"
printf '  %q\n' "${CMD[@]}"

# Run the node
exec "${CMD[@]}"
