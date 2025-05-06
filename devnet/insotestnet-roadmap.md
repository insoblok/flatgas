
# Insotestnet Roadmap

## Phase 0: Bootstrapping Local Testnet (Insotestnet)

### 🎯 Goals
- Simulate a realistic multi-node devnet.
- Establish `inso` binary configuration norms.
- Verify basic peer connectivity and transaction propagation.

---

## ✅ Step 1: Initialize Genesis Config

- Define `genesis.json` with required chainID, allocations, and config.
- Copy it to all nodes.
- Initialize all nodes using `inso init genesis.json`.

---

## ✅ Step 2: Configure Static/Boot Nodes

- Option A: Use `--bootnodes` CLI flag.
- Option B: Use `static-nodes.json` file.
- Option C: Use discovery and rely on internal DNS (Docker).

Recommended for testnet: **Option B**.

---

## ✅ Step 3: Create Scripts

- `entrypoint-node1.sh`: initialize and start with mining.
- `entrypoint-node2.sh`: connect to node1 using bootnode.
- `attach-to-nodeX.sh`: attach Geth JS console.

---

## ✅ Step 4: Docker Compose Setup

- `docker-compose.yml` to manage node containers.
- Mount volumes for keys, data, and genesis.
- Expose ports for RPC/WebSocket/P2P.

---

## ✅ Step 5: Peer Verification

- Use `admin.peers` on node consoles to verify P2P connections.
- Test transactions and logs.

---

## ✅ Step 6: Docs & Maintenance

- Document config, genesis file, and node setup process.
- Version the `genesis.json` and `static-nodes.json`.

---

## 🔁 Step 7: Iteration and Scaling

- Add node3+ using `static-nodes.json`.
- Observe peer discovery effects.
- Test network behavior with one node down.

---

## ⛳️ Migration to Production

- Fork from tested genesis.
- Adjust validator list and initial allocation.
- Remove dev-only accounts.
- Set stricter configs and permissions.

---

## Notes

- Enforce consistent `networkid` across nodes.
- Maintain enode URLs and nodekeys per environment.
