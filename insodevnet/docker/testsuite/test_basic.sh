#!/bin/sh
set -e

echo "🔍 Running Flatgas Basic Node Test Suite"

# Correct IPC path in container
IPC="/app/data/geth.ipc"

# Function to run JavaScript on the validator node
run_inso() {
  docker exec flatgas-validator1 /usr/local/bin/inso attach --exec "$1" "$IPC"
}

echo "🔢 Block Number:"
run_inso "eth.blockNumber"

echo "👥 Peer Count:"
run_inso "net.peerCount"

echo "🗳️ Clique Signers:"
run_inso "clique.getSigners()"
