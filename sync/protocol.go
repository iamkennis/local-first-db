package sync


import (
	"encoding/json"

	"github.com/iamkennis/decentralized-db/core"
)

func Encode(op core.Operation) ([]byte, error) {
	return json.Marshal(op)
}

func Decode(data []byte) (core.Operation, error) {
	var op core.Operation
	err := json.Unmarshal(data, &op)
	return op, err
}