# Directory: devnet/

# --- genesis.json ---
# Save as: devnet/genesis.json
{
  "config": {
    "chainId": 12345,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "berlinBlock": 0,
    "londonBlock": 0,
    "clique": {
      "period": 5,
      "epoch": 30000
    }
  },
  "difficulty": "1",
  "gasLimit": "8000000",
  "alloc": {
    "0x34b3248010CcCf160838Ea336D598b8747ddd147": {
      "balance": "1000000000000000000000000"
    }
  },
  "extraData": "0x000000000000000000000000000000000000000000000000000000000000000034b3248010cccf160838ea336d598b8747ddd14700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "nonce": "0x0",
  "timestamp": "0x0",
  "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "coinbase": "0x0000000000000000000000000000000000000000",
  "number": "0x0",
  "gasUsed": "0x0",
  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
}

# --- init-chain.sh ---
# Save as: devnet/scripts/init-chain.sh
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

# --- start-devnet.sh ---
# Save as: devnet/scripts/start-devnet.sh
#!/bin/bash
set -e
exec ./scripts/run-devnet.sh ./devnet/data ./build/bin/inso

# --- create-account.sh ---
# Save as: devnet/scripts/create-account.sh
#!/bin/bash
set -euo pipefail

INSO=${1:-"../../build/bin/inso"}
KEYDIR=${2:-"../keys"}

mkdir -p "$KEYDIR"

echo "🔐 Creating new account in keystore: $KEYDIR"
echo "🧪 Command: $INSO account new --keystore $KEYDIR"

exec "$INSO" account new --keystore "$KEYDIR"

# --- make-extradata.sh ---
# Save as: devnet/scripts/make-extradata.sh
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
echo "🤔 extraData for signer:"
echo "$EXTRADATA"

# --- .gitignore ---
# Save as: devnet/.gitignore
/data/
/keys/

# --- README.md ---
# Save as: devnet/README.md

# Flatgas Devnet

This directory contains everything needed to bootstrap and run a local Flatgas devnet node.

## Quickstart

1. **Create a validator account:**
   ```bash
   ./devnet/scripts/create-account.sh
   ```
   > Save the generated address and keep the key safe!

2. **Generate `extraData` field for your address:**
   ```bash
   ./devnet/scripts/make-extradata.sh 0x<your-signer-address>
   ```
   > Copy the output and replace the `extraData` field in `genesis.json`

3. **Initialize the chain:**
   ```bash
   ./devnet/scripts/init-chain.sh
   ```

4. **Start the node:**
   ```bash
   ./devnet/scripts/start-devnet.sh
   ```

The node will launch with Clique PoA, custom chain ID 12345, and full RPC support.

> You can use this as a real testbed to simulate your mainnet configuration.
