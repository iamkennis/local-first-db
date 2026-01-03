package main

import (
	"fmt"
	"log"

	"github.com/iamkennis/decentralized-db/core"
	"github.com/iamkennis/decentralized-db/storage/file"
)

func main() {
	fmt.Println("🚀 Decentralized DB - Local Test App")

	// Create file-based storage
	storage := file.New("local_test.log")

	// Create store
	store, err := core.NewStore(storage)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	// Example operations
	store.Apply(core.Operation{
		ID:        "1",
		Actor:     "local-device",
		Timestamp: 1,
		Type:      "set",
		Key:       "welcome",
		Value:     []byte("Hello from restructured project!"),
	})

	// Get value
	val := store.Get("welcome")
	fmt.Printf("✅ Retrieved: %s\n", string(val))
	fmt.Println("📊 Store initialized successfully")
}
