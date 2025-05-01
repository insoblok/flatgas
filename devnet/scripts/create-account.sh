#!/bin/bash
set -euo pipefail

INSO=${1:-"../../build/bin/inso"}
KEYDIR=${2:-"../keys"}

mkdir -p "$KEYDIR"

echo "🔐 Creating new account in keystore: $KEYDIR"
echo "🧪 Command: $INSO account new --keystore $KEYDIR"

exec "$INSO" account new --keystore "$KEYDIR"
