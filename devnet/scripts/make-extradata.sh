#!/bin/bash
set -euo pipefail

ADDRESS=${1:-}

if [[ -z "$ADDRESS" ]]; then
  echo "Usage: $0 <signer-address>"
  echo "Example: $0 0x34b3248010CcCf160838Ea336D598b8747ddd147"
  exit 1
fi

ADDR_HEX=$(echo "$ADDRESS" | sed 's/^0x//' | tr '[:upper:]' '[:lower:]')

if [[ ${#ADDR_HEX} -ne 40 ]]; then
  echo "Error: Address must be 40 hex chars (20 bytes)"
  exit 1
fi

VANITY=$(printf '00%.0s' {1..64})
SEAL=$(printf '00%.0s' {1..130})

EXTRADATA="0x${VANITY}${ADDR_HEX}${SEAL}"
echo "🧐 extraData for signer:"
echo "$EXTRADATA"
