#!/bin/bash
set -euo pipefail

INSO=$(realpath "${1:-../../build/bin/inso}")
DATADIR=$(realpath "${2:-../data}")
GENESIS=$(realpath "${3:-../genesis.json}")

echo "🧱 Initializing devnet chain"
echo "🔧 Binary:   $INSO"
echo "📁 Datadir:  $DATADIR"
echo "📄 Genesis:  $GENESIS"
echo "🧪 Command:"
echo "  $INSO --datadir $DATADIR init $GENESIS"

exec "$INSO" --datadir "$DATADIR" init "$GENESIS"
