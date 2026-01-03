//go:build js && wasm
package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/iamkennis/decentralized-db/core"
)

var store = core.NewStore()

func applyOp(this js.Value, args []js.Value) any {
	var op core.Operation
	json.Unmarshal([]byte(args[0].String()), &op)
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
