# 🐳 Flatgas Devnet — Docker Setup

This guide describes how to run the Flatgas devnet using Docker.

---

## 📁 Directory Layout

Ensure this structure exists inside the Flatgas repo:

```
flatgas/
├── build/bin/inso             # Compiled binary
├── devnet/
│   ├── Dockerfile             # Node image builder
│   ├── docker-compose.yml     # Multi-node local testnet
│   ├── genesis.json
│   ├── keys/                  # Validator keys
│   ├── data/                  # Chain data (auto-generated)
│   └── scripts/
│       └── build-docker.sh    # Image builder script
```

---

## 🛠️ Build the Docker Image

Run this from the **repo root**:

```bash
./devnet/scripts/build-docker.sh
```

This builds the image using the context from the root so `build/bin/inso` is available.

---

## 🧪 Next Steps

After building, continue with:

- `docker-compose.yml` setup
- Running multiple nodes
- Connecting external peers or clients

---

## ⚠️ Notes

- Always run Docker-related scripts from the root (`flatgas/`)
- Do **not** run from inside `devnet/scripts/` — relative paths will break
