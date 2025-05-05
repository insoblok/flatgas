#!/bin/sh
set -e

if [ ! -d "/app/data/geth/chaindata" ]; then
  echo "🧱 Initializing genesis block..."
  inso --nodekey /app/keys/nodekey --datadir /app/data init /app/genesis.json
fi

echo "🚀 Starting node1..."
exec inso --nodekey /app/keys/nodekey --datadir /app/data \
  --http --http.api admin,debug,eth,miner,net,txpool,web3 \
  --http.addr 0.0.0.0 \
  --http.port 8545 \
  --http.corsdomain "*" \
  --ws --ws.api eth,net,web3 \
  --ws.addr 0.0.0.0 \
  --ws.port 8546 \
  --ws.origins "*" \
  --mine \
  --nodiscover \
  --networkid 12345 \
  --allow-insecure-unlock
