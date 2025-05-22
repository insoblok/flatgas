# Flatgas `txsender` CLI Tool – Roadmap

This roadmap tracks the development of the `txsender` CLI tool for managing Flatgas devnet accounts, transactions, and node interaction.

---

## 🧱 CLI Command Plan

### 1. `config`
- Load/save CLI configuration (`wallet/config.json`)
- Define named RPCs (e.g. `local`, `devnet`)
- Set defaults (e.g. `defaultRpc`, `defaultFrom`)
- Commands:
    - `txsender config add-rpc --name devnet --url http://...`
    - `txsender config set-default-rpc devnet`
    - `txsender config list-rpcs`

### 2. `ping`
- Connect to RPC
- Show latest block number and basic info

### 3. `new`
- Create a new account in the keystore
- Store alias → address in `aliases.json`

### 4. `fund`
- Send ETH from faucet account to a recipient (alias or address)

### 5. `list`
- List all known accounts with aliases and addresses

### 6. `balance`
- Check ETH balance for an alias or address

### 7. `send`
- Send ETH from one alias/address to another

---

## 📁 CLI Structure

```bash
txsender/
├── main.go
├── cmd/
│   ├── config.go
│   ├── ping.go
│   ├── new.go
│   ├── fund.go
│   ├── list.go
│   ├── balance.go
│   └── send.go
├── wallet/
│   ├── config.json
│   ├── aliases.json
│   └── keystore/
├── internal/
│   ├── accounts.go
│   ├── config.go
│   ├── rpc.go
│   └── tx.go
```

---

## 🧩 Optional Future Features

- `import` / `export` of keys
- `menu` mode for guided interaction
- HTTP faucet server
