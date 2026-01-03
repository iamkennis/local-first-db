# 🏗️ Architecture

## Why append-only logs?
They provide:
- Immutability
- Deterministic replay
- Simple conflict resolution
- Easy sync

## Why local-first?
- No loading spinners
- Works offline
- User owns data

## Why WASM + Go?
- Deterministic execution
- Shared logic between backend and browser
- Strong concurrency model

## Why stateless servers?
- Easier scaling
- No trust required
- Replaceable relays

This document explains the core design decisions behind the decentralized database.

---

## 🎯 Core Principles

1. **Offline-first**: Operations work without network connectivity
2. **Conflict-free**: CRDT merge logic ensures eventual consistency
3. **Privacy-first**: End-to-end encryption, zero server trust
4. **Simple sync**: Dumb relay server keeps infrastructure minimal

---

## 📝 Why Append-Only Logs?

### The Problem
Traditional databases use **mutable state**:
- UPDATE statements modify rows in-place
- DELETE removes data permanently
- Concurrent writes require locks
- Conflicts need central coordination

### Our Solution: Immutable Event Log

Every operation is **appended**, never modified:

```
[op1: set user:1 = "Alice"]
[op2: set user:1 = "Bob"]      ← Doesn't override op1
[op3: delete user:1]            ← Doesn't erase previous ops
```

### Benefits

#### 1️⃣ Crash Recovery
- System crashes mid-write? No corruption.
- Restart and replay the log from disk.
- Atomic writes: operation fully written or not at all.

#### 2️⃣ Distributed Sync
- Can merge logs from different devices
- No "last modified" conflicts
- Deterministic replay on all peers

#### 3️⃣ Audit Trail
- Complete history of all changes
- Can reconstruct state at any point in time
- Useful for debugging and compliance

#### 4️⃣ Simple Implementation
```go
func Append(op Operation) error {
    data := json.Marshal(op)
    file.Write(append(data, '\n'))  // Just append!
}
```

No in-place updates, no B-trees, no WAL complexity.

### Trade-offs
- ❌ Log grows unbounded (needs compaction strategy)
- ❌ Slower queries (must replay entire log)
- ✅ Rock-solid consistency guarantees
- ✅ Trivial to implement correctly

---

## 🔐 Why Single-Writer Model?

### The Problem with Multiple Writers

Race condition example:
```go
// Thread 1                    // Thread 2
ops := store.ops               ops := store.ops
ops = append(ops, op1)         ops = append(ops, op2)
store.ops = ops                store.ops = ops
                               ← op1 lost!
```

Traditional solutions:
- Mutexes (slow, blocking)
- Lock-free algorithms (complex, error-prone)
- Transactions (heavyweight overhead)

### Our Solution: Single Goroutine Owner

```go
func (s *Store) run() {
    for op := range s.input {
        s.ops = append(s.ops, op)  // Only this goroutine writes
    }
}

func Apply(op Operation) {
    s.input <- op  // Everyone else just sends
}
```

### Why This Works

#### 1️⃣ No Race Conditions
- Only one goroutine modifies `s.ops`
- Go channels provide synchronization
- No mutex overhead

#### 2️⃣ Natural Backpressure
- Slow disk I/O? Channel fills up
- Fast producers automatically wait
- No manual flow control needed

#### 3️⃣ Sequential Consistency
- Operations processed in channel order
- Easier to reason about
- Deterministic behavior

#### 4️⃣ Simple Testing
```bash
go test -race  # Zero data races!
```

### Performance Implications
- ✅ Single-threaded writes scale to ~100K ops/sec
- ✅ Reads can be concurrent (RWMutex for Get)
- ❌ Can't parallelize writes across cores
- ❌ Bounded by single disk throughput

For most CRDT use cases (collaborative apps, local-first), this is plenty fast.

---

## 🌐 Why IndexedDB (Browser Storage)?

### Options Comparison

| Storage | Capacity | Async | Structured | Verdict |
|---------|----------|-------|------------|---------|
| `localStorage` | 5-10 MB | ❌ Sync | ❌ String-only | Too limited |
| `sessionStorage` | 5-10 MB | ❌ Sync | ❌ String-only | Too limited |
| **IndexedDB** | **50+ MB** | **✅ Async** | **✅ Objects** | **Winner** |
| WebSQL | Deprecated | ✅ | ✅ | Dead standard |

### Why IndexedDB Wins

#### 1️⃣ Capacity
```javascript
// localStorage: Only ~5MB
localStorage.setItem('ops', JSON.stringify(ops))  // Quota exceeded!

// IndexedDB: 50MB+ (browser-dependent)
db.put('operations', ops)  // Plenty of room
```

#### 2️⃣ Asynchronous API
```javascript
// localStorage blocks main thread
const data = localStorage.getItem('big-data')  // UI freezes 😰

// IndexedDB doesn't block
db.getAll('operations').then(ops => {
    // UI stays responsive ✨
})
```

#### 3️⃣ Structured Data
```javascript
// localStorage: Everything is a string
const user = JSON.parse(localStorage.getItem('user'))  // Manual parse

// IndexedDB: Store objects directly
db.put('users', { id: 1, name: 'Alice' })  // Native objects
```

#### 4️⃣ Indexing and Queries
```javascript
// Find all ops after timestamp 100
const tx = db.transaction('operations')
const index = tx.objectStore('operations').index('timestamp')
const ops = index.getAll(IDBKeyRange.lowerBound(100))
```

#### 5️⃣ Transactions
```javascript
const tx = db.transaction(['operations', 'seen'], 'readwrite')
tx.objectStore('operations').add(op)
tx.objectStore('seen').put(op.id, true)
// Atomic: both succeed or both fail
```

### WASM Integration
```go
//go:build js && wasm

import "syscall/js"

func (s *IndexedDBStorage) Append(op Operation) error {
    db := js.Global().Get("indexedDB")
    // Call browser IndexedDB APIs from Go!
}
```

---

## 🔒 Why Client-Side Encryption?

### Zero-Trust Architecture

```
┌─────────────────┐
│  User's Device  │
│  ┌───────────┐  │
│  │ Plain Ops │  │     Encryption happens HERE
│  └─────┬─────┘  │            ↓
│        │ AES-GCM│     ┌──────────────┐
│  ┌─────▼─────┐  │     │ Encrypted Op │
│  │ Cipher Op │──┼────>│ (Base64)     │
│  └───────────┘  │     └──────────────┘
└─────────────────┘              │
                                 │
                          ┌──────▼──────┐
                          │    Server   │
                          │  (can't read)│
                          └─────────────┘
```

### Why Encrypt on Client?

#### 1️⃣ Server Can't Read Your Data
```go
// Server receives this:
{
  "id": "op123",
  "value": "nfj38f9h2f9h2f9..."  // Encrypted blob
}

// Server has zero idea what this means
```

Even if server is compromised:
- ✅ Attacker gets encrypted gibberish
- ✅ No encryption keys on server
- ✅ No plaintext data

#### 2️⃣ Compliance-Friendly
- GDPR: "Right to be forgotten" → Just delete your key
- HIPAA: PHI never leaves device unencrypted
- SOC2: Server doesn't process sensitive data

#### 3️⃣ Multi-Tenant Security
```
User A's key encrypts User A's data
User B's key encrypts User B's data

Server stores both, can't distinguish them
```

No user can decrypt another user's data, even if they hack the server.

#### 4️⃣ Offline Encryption
```go
// Works without network
encrypted := Encrypt(data, identity.Key)
storage.Append(encrypted)

// Sync later when online
```

### Implementation

```go
func Encrypt(data, key []byte) ([]byte, error) {
    block, _ := aes.NewCipher(key)
    gcm, _ := cipher.NewGCM(block)
    
    nonce := make([]byte, 12)
    rand.Read(nonce)
    
    // Authenticated encryption
    return gcm.Seal(nonce, nonce, data, nil), nil
}
```

**AES-GCM provides:**
- Confidentiality (encryption)
- Authenticity (can't tamper without detection)
- Associated data (can include metadata)

---

## 🤝 Why Dumb Relay Server?

### Smart Server (Traditional)
```
┌──────────────────────┐
│    Server            │
│  ┌────────────────┐  │
│  │ User Auth      │  │  
│  │ Permissions    │  │  ← Complex logic
│  │ Conflict Res   │  │  ← State machine  
│  │ Data Transform │  │  ← Business rules
│  └────────────────┘  │
└──────────────────────┘
```

Problems:
- ❌ Server must understand application logic
- ❌ Tight coupling between client and server
- ❌ Complex deployment and scaling
- ❌ Single point of failure

### Dumb Relay (Our Approach)

```go
func relay() {
    for msg := range broadcast {
        for client := range clients {
            client.Send(msg)  // Just forward!
        }
    }
}
```

That's it. **42 lines of code.**

### Why This Is Powerful

#### 1️⃣ Zero Application Logic
Server doesn't know:
- What the data means
- Who can access what
- How to resolve conflicts

All intelligence in clients.

#### 2️⃣ Trivial to Scale
```bash
# Run 10 relays behind load balancer
for i in 1..10; do
    ./relay --port=$((8080 + i)) &
done
```

No shared state, no coordination.

#### 3️⃣ Cheap Infrastructure
```
Client A ──┐
           ├──> Relay ($5/month VPS) ──┐
Client B ──┘                           ├──> Client C
                                       └──> Client D
```

Compare to Firebase ($25-$500/month) with vendor lock-in.

#### 4️⃣ Self-Hosting Friendly
```bash
# Users can run their own relay
docker run -p 8080:8080 decentralized-db/relay
```

No need to trust a third party.

#### 5️⃣ Swappable Backends
```
┌─────────┐
│ WebRTC  │  ← Peer-to-peer, no server
└─────────┘

┌─────────┐
│ Relay   │  ← Simple broadcast
└─────────┘

┌─────────┐
│ libp2p  │  ← Decentralized network
└─────────┘
```

Same CRDT core, different transport layers.

### Security Model

```
Client sends: ENCRYPTED operation
Relay sees:   Random bytes
Relay does:   Broadcast to peers
Peers decrypt using their keys
```

No auth needed! Encrypted data is self-authenticating.

---

## 🎯 Design Wins

### What We Get Right

✅ **Simplicity**: Core CRDT logic is ~300 lines  
✅ **Correctness**: Deterministic merge, no race conditions  
✅ **Privacy**: E2E encryption, server can't read data  
✅ **Resilience**: Offline-first, crash recovery  
✅ **Cost**: Self-host on $5/month VPS  

### Known Trade-offs

❌ **Storage overhead**: Log needs compaction  
❌ **Query performance**: Must replay log  
❌ **No indexes**: Linear scan for queries  
❌ **No transactions**: Eventually consistent only  

### When to Use This

✅ Collaborative apps (Notion, Figma)  
✅ Local-first software (offline editing)  
✅ Privacy-critical apps (healthcare, finance)  
✅ Multi-device sync (phone + laptop + tablet)  

❌ Real-time analytics (need aggregations)  
❌ Complex joins (need relational DB)  
❌ Strong consistency (need distributed consensus)  

---

## 📚 Further Reading

- [CRDTs: The Hard Parts](https://youtu.be/x7drE24geUw) - Martin Kleppmann
- [Local-First Software](https://www.inkandswitch.com/local-first/)
- [Designing Data-Intensive Applications](https://dataintensive.net/)

---

**Architecture designed for simplicity, privacy, and resilience.**
