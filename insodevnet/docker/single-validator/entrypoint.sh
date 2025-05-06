#!/bin/sh
set -e

if [ ! -d "/app/data/geth/chaindata" ]; then
  echo "🧱 Initializing genesis..."
  inso --datadir /app/data init /app/genesis.json
fi

echo "🚀 Starting validator1 node..."
exec inso --datadir /app/data \
  --nodekey /app/nodekey/nodekey \
  --networkid 12345 \
  --http --http.addr 0.0.0.0 --http.port 8545 \
  --http.api eth,net,web3,admin,clique \
  --nodiscover \
