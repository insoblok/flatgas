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

## 🔜 Planned Features

### 🔄 Transactions
- [ ] `tx receipt` — detailed transaction result, logs, gas used
- [ ] `tx watch` — polling-based monitor until TX is mined (in `txsender`)

### 🌐 Node and Chain
- [ ] `block latest` — get latest block metadata
- [ ] `node peers` — list active P2P peers

### 🔐 Wallet Enhancements
- [ ] `accounts delete` — remove alias and keystore
- [ ] `accounts rename` — rename alias

## 📡 WebSocket & Event Monitoring (Future)

Planned for separate tool/module: `insowatch`
- Watch contract logs, pending TXs, block headers
- Uses `eth_subscribe` via WebSocket
- Can emit logs, webhooks, or write to file

---

The core CLI is complete and stable. Next steps focus on network awareness, event monitoring, and usability improvements.
