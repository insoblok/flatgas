# Flatgas 🛠️ txsender CLI Roadmap

This roadmap defines the command-line interface structure and functionality for managing Flatgas devnet wallets and transactions.

## ✅ Completed Commands

### 📁 Config
- [x] `config list-rpcs` — show configured RPC endpoints
- [x] `config add-rpc` — add or update an RPC alias
- [x] `config set-default-rpc` — define which RPC to use by default

### 🔐 Accounts
- [x] `accounts create` — create new keystore + alias
- [x] `accounts list` — list known accounts by alias
- [x] `accounts balance` — check ETH balance by alias or address
- [x] `accounts import` — import existing keystore JSON and alias
- [x] `accounts export` — export keystore by alias

### 💸 Transactions
- [x] `fund` — send ETH from known account to another alias/address
- [x] `fund` (with `--send`) — performs real transfer
- [x] `fund` (default dry-run) — safety mechanism

### 🔎 Transaction Status
- [x] `tx status` — check if a transaction has been mined and its details

### 🖥️ Node Interaction
- [x] `node info` — show current chain ID, latest block, gas price, peer count

## 🔜 Optional / Future

### ⛓️ Chain Monitoring
- [ ] `tx watch` — live poll a TX until mined
- [ ] `tx receipt` — extended version of `tx status`

### 🧑‍💻 UX Improvements
- [ ] `menu` — TUI-based interactive interface (optional)
- [ ] `accounts delete` / `rename`

### 🌐 Optional Network Utilities
- [ ] `node peers` — show connected P2P peers
- [ ] `block latest` — fetch latest block header or number

---

The core CLI is complete and testnet-ready. 🎉
