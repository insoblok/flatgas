# InsoDevnet Roadmap

This document tracks the plan and evolution of the `insodevnet`, a Flatgas-powered test network intended to explore validator setup, network rules, and production-readiness.


---

## 🚀 Roadmap

### 🧱 Phase 1 — Local Single-Validator Clique Chain

* [x] Initialize `genesis.json` with Clique consensus
* [x] Run 1 node with Docker (`inso`) using `--mine`
* [x] Validate genesis inclusion, block production
* [x] CLI test: deploy contract, send tx, check balance

### ☁️ Phase 1B — Cloud Single-Validator

* [ ] Deploy same Docker setup on cloud VM
* [ ] Setup remote access to RPC (with auth if needed)
* [ ] Allow team members to connect wallets

### 🧱 Phase 2 — Local Multi-Validator Clique Chain

* [ ] Add multiple nodes with static validator keys in genesis
* [ ] Configure shared network ID, enode bootstrapping
* [ ] Validate consensus rotation by block number

### ☁️ Phase 2B — Cloud Multi-Validator

* [ ] Docker-compose setup across multiple cloud VMs
* [ ] Enable metrics/logging per validator
* [ ] Optional: Add Grafana/Prometheus later

### 🔁 Phase 3 — Experiment with Validator Dynamics

* [ ] Add validator via Clique vote
* [ ] Remove validator
* [ ] Observe network convergence, block production

### 🔀 Phase 4 — Switch to PoS (if needed)

* [ ] Prepare new PoS-compatible genesis
* [ ] Implement stake/validator set logic (if supported)
* [ ] Document differences from Clique

---

## 🔄 Migration Plan: Testnet → Mainnet

* Keep all chain rules, chain ID, validator design modular
* Use separate genesis, keys, and chain ID for `insomainnet`
* Use same dockerized tooling for deployment
* Bake in upgradeability, if governance or config voting required

---

## 🧠 Notes & Best Practices

* Always separate data dirs per node
* Document genesis + key origins clearly
* Prefer `static-nodes.json` in controlled setups
* Centralize secrets if deploying in cloud (e.g., Vault)
* Use tmux/screen when SSH-ing into nodes for long sessions

---

*Maintained by: Flatgas Core / Iyad — 2025*
