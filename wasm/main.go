//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/iamkennis/decentralized-db/core"
)

var store = core.NewStore()
var userKey = make([]byte, 32) // Will be set from JS

func applyOp(this js.Value, args []js.Value) any {
	encrypted := make([]byte, args[0].Length())
	js.CopyBytesToGo(encrypted, args[0])

	plain, err := core.Decrypt(userKey, encrypted)
	if err != nil {
		return nil
	}

	var op core.Operation
	json.Unmarshal(plain, &op)

	store.Apply(op)
	return nil
}

func state(this js.Value, args []js.Value) any {
	b, _ := json.Marshal(store.State())
	return string(b)
}

func main() {
	js.Global().Set("applyOp", js.FuncOf(applyOp))
	js.Global().Set("getState", js.FuncOf(state))
	select {}
}
