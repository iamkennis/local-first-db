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

Browser
┌──────────────────────────────┐
│ UI (JS / React)              │
│                              │
│ IndexedDB (append-only log)  │
│        ▲                     │
│        │                     │
│ Go Core (WASM)               │
│ - Ops                        │
│ - Merge                      │
│ - Crypto                     │
└────────┬─────────────────────┘
│
▼
Dumb WebSocket Relay

Browser
┌──────────────────────────────┐
│ UI (JS / React)              │
│                              │
│ IndexedDB (append-only log)  │
│        ▲                     │
│        │                     │
│ Go Core (WASM)               │
│ - Ops                        │
│ - Merge                      │
│ - Crypto                     │
└────────┬─────────────────────┘
│
▼
Dumb WebSocket Relay

core/        → Pure database logic (Go)
storage/     → File + IndexedDB backends
sync/        → WebSocket protocol
wasm/        → Go → WASM bridge
cmd/relay/   → Stateless relay server

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

