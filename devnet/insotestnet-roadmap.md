# Flatgas `insotestnet` Roadmap

## 🔮 Overview

This document defines the setup plan for launching and operating `insotestnet`, a Flatgas-based internal test network. It also outlines the steps to migrate cleanly to a production-grade mainnet (`inso-mainnet`).

---

## 🔢 Phase 1: Bootstrap `insotestnet`

### 1. Prepare `genesis.json`

* Define chain ID (e.g., `12345`)
* Configure block time, validator list, premine
* Save at: `flatgas/networks/insotestnet/genesis.json`

### 2. Define Validator Set

* Generate `nodekey` files for each validator
* Store in `flatgas/networks/insotestnet/keys/`
* Build `static-nodes.json` listing peer enodes

### 3. Docker Compose Per Node

* Compose files: `docker-compose.node1.yml`, etc.
* Wait for `node1` to fully start, then bring up other nodes

### 4. Bootstrap Scripts

* Create `scripts/start-insotestnet.sh`
* Automate:

    * Chain initialization (`inso init`)
    * Static peer config
    * Key mounting

### 5. Observe the Network

* Use Geth console `admin.peers`, logs
* Optional: run lightweight explorer or Prometheus/Grafana later

### 6. Simulate Governance

* Dummy proposals for flat fee review
* Track epochs, simulate changes manually or with CLI tools

### 7. Deploy Sample Transactions

* Send simple txs
* Optionally deploy test contracts (if supported)
* Observe mempool, inclusion, confirmation latency

---

## 🚀 Phase 2: Transition to Production (`inso-mainnet`)

### What Can Be Reused

* `genesis.json` format and logic
* Validator key generation + mounting system
* Docker Compose and init scripts
* Monitoring tools

### What Must Be Changed

* Harden `genesis.json`: no premine, updated validator list
* Key handling: use secrets, HSMs or encrypted stores
* Make `static-nodes.json` public or exposed via DNS seed
* Add RPC protection (optional): rate limiting, TLS
* Publish validator onboarding guides

### Launch Checklist

*

---

## ✅ Best Practices Summary

| Goal                | Best Practice                          |
| ------------------- | -------------------------------------- |
| Bootstrap stability | Use `static-nodes.json`, no discovery  |
| Peer hygiene        | Start with few nodes, expand gradually |
| Key management      | Automate or securely store nodekeys    |
| Reproducibility     | Use scripts + tagged Docker images     |
| Upgrade handling    | Plan versioning + hard fork activation |

---

## 📂 Suggested Layout

```
flatgas/
├── networks/
│   └── insotestnet/
│       ├── genesis.json
│       ├── static-nodes.json
│       ├── keys/
│       │   ├── nodekey1
│       │   └── nodekey2
│       └── compose/
│           ├── docker-compose.node1.yml
│           └── docker-compose.node2.yml
├── scripts/
│   └── start-insotestnet.sh
```
