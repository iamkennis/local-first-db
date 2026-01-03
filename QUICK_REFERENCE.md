# 📁 Quick Reference - New Structure

## Import Paths

```go
// Core CRDT logic
import "github.com/iamkennis/decentralize-db/core"

// File storage
import "github.com/iamkennis/decentralize-db/storage/file"

// Future: Sync protocol  
import "github.com/iamkennis/decentralize-db/sync"
```

## Usage Example

```go
package main

import (
    "github.com/iamkennis/decentralize-db/core"
    "github.com/iamkennis/decentralize-db/storage/file"
)

func main() {
    // Create storage backend
    storage := file.NewFileStorage("my_database.log")
    
    // Create store
    store, _ := core.NewStore(storage)
    
    // Apply operations
    store.Apply(core.Operation{
        ID:        "op1",
        Actor:     "device-1",
        Timestamp: 123,
        Type:      "set",
        Key:       "user:name",
        Value:     []byte("Alice"),
    })
    
    // Retrieve data
    value := store.Get("user:name")
}
```

## Directory Map

| Path | Purpose | Status |
|------|---------|--------|
| `core/` | Pure CRDT logic | ✅ Done |
| `storage/file/` | File persistence | ✅ Done |
| `storage/indexeddb/` | Browser storage | 🚧 Phase 6 |
| `sync/` | Network sync | 🚧 Phase 4 |
| `wasm/` | Browser build | 🚧 Phase 6 |
| `cmd/local/` | Test app | ✅ Working |
| `cmd/relay/` | Relay server | 🚧 Phase 4 |

## Running Apps

```bash
# Local test app (works now!)
go run cmd/local/main.go

# Relay server (Phase 4)
go run cmd/relay/main.go
```

## Old Phase-1 Location

All original work preserved in:
```
phase-1/
├── All .go files
├── Tests  
├── Documentation (.md)
└── Experiment guides
```

Next task: Migrate tests to new structure!
