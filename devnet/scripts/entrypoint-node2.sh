#!/bin/sh
set -e

echo "⛓️ Connecting node2 to bootstrap node..."

exec inso --datadir /app/data \
  --http --http.api admin,debug,eth,net,txpool,web3 \
  --http.addr 0.0.0.0 \
  --http.port 8545 \
  --http.corsdomain "*" \
  --ws --ws.api eth,net,web3 \
  --ws.addr 0.0.0.0 \
  --ws.port 8546 \
  --ws.origins "*" \
  --bootnodes enode://6699c5a736e5f9a0699b0c10f94bbd564b32a7c84518b7b1d355fbc96fd4fa2116a88c9a30c4e4359f3ea13950b29e843bfa47d8aa921daef025ac0654398a02@flatgas-node1:30303 \
  --networkid 12345 \
  --allow-insecure-unlock
