# Local-First Database (Go + WASM + IndexedDB)

A **local-first, offline-first database** built from scratch using **Go**, **WebAssembly**, **IndexedDB**, and an **append-only operation log**.

The system is designed to be:
- Deterministic
- Eventually consistent
- Zero-knowledge (server never sees plaintext)
- Decentralized by architecture, not hype

---

## Why this project?

Most apps still rely on:
- Centralized databases
- Constant network access
- Server-authoritative state

This project explores an alternative model:

> **State is derived. Operations are the source of truth.**

Inspired by systems like Figma, Notion, Linear, and modern local-first architectures.

---

## Core Concepts

- **Event sourcing** (append-only log)
- **Deterministic conflict resolution**
- **Single-writer concurrency model**
- **Client-side encryption**
- **Stateless sync servers**
- **Local-first browser storage**

---

## Architecture Overview

## Architecture Overview

```
┌─────────────────────┐                    ┌─────────────────────┐
│   Browser A         │                    │   Browser B         │
│  ┌───────────────┐  │                    │  ┌───────────────┐  │
│  │  UI (React)   │  │                    │  │  UI (React)   │  │
│  └───────┬───────┘  │                    │  └───────┬───────┘  │
│          │          │                    │          │          │
│  ┌───────▼───────┐  │                    │  ┌───────▼───────┐  │
│  │   Go (WASM)   │  │                    │  │   Go (WASM)   │  │
│  │ ┌───────────┐ │  │                    │  │ ┌───────────┐ │  │
│  │ │ Store     │ │  │                    │  │ │ Store     │ │  │
│  │ │ Merge     │ │  │                    │  │ │ Merge     │ │  │
│  │ │ Crypto    │ │  │                    │  │ │ Crypto    │ │  │
│  │ └───────────┘ │  │                    │  │ └───────────┘ │  │
│  └───────┬───────┘  │                    │  └───────┬───────┘  │
│          │          │                    │          │          │
│  ┌───────▼───────┐  │                    │  ┌───────▼───────┐  │
│  │  IndexedDB    │  │                    │  │  IndexedDB    │  │
│  │ (append-only) │  │                    │  │ (append-only) │  │
│  └───────────────┘  │                    │  └───────────────┘  │
└──────────┬──────────┘                    └──────────┬──────────┘
           │                                          │
           │ Encrypted ops                            │ Encrypted ops
           │ (WebSocket)                              │ (WebSocket)
           │                                          │
           └──────────────┬───────────────────────────┘
                          │
                          ▼
                 ┌────────────────┐
                 │  Relay Server  │
                 │   (Stateless)  │
                 │                │
                 │  • Broadcast   │
                 │  • No storage  │
                 │  • No auth     │
                 └────────────────┘
```

**Key Points:**
- Each client has its own local IndexedDB storage
- Go core compiled to WASM runs in browser
- Relay server just forwards encrypted operations
- No server-side logic or state

### Project Structure

```
decentralized-db/
├── core/                    # Pure CRDT logic (Go)
│   ├── operation.go         # Operation definition
│   ├── merge.go             # Deterministic conflict resolution
│   ├── store.go             # Single-writer store
│   ├── storage.go           # Storage interface
│   ├── identity.go          # Device identity + keys
│   └── crypto.go            # AES-GCM encryption
│
├── storage/
│   ├── file/                # File-based persistence
│   │   └── file_storage.go  # Append-only log
│   └── indexeddb/           # Browser storage (WASM)
│       └── indexeddb.js     # JS bridge to IndexedDB
│
├── sync/                    # Network sync
│   ├── protocol.go          # Operation encoding
│   └── client.go            # WebSocket client
│
├── wasm/                    # Browser build
│   ├── main.go              # Go → WASM entry point
│   └── index.html           # WASM loader
│
└── cmd/
    ├── relay/               # WebSocket relay server
    │   └── main.go          # Stateless broadcast
    └── local/               # CLI demo app
        └── main.go          # Test operations
```

---

## How to Run

### 1. Run relay server
```bash
go run cmd/relay/main.go
```

### 2. Build WASM
```bash
GOOS=js GOARCH=wasm go build -o wasm/main.wasm ./wasm
```

### 3. Serve WASM
```bash
serve ./wasm
```

