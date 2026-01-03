package core

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type Operation struct {
	ID        string
	Key       string
	Value     string
	Timestamp int64
	ActorID   string
}

func RandID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func NewOp(key, value string) Operation {
	return Operation{
		ID:        RandID(),
		Key:       key,
		Value:     value,
		Timestamp: time.Now().UnixNano(),
	}
}

func newer(a, b Operation) bool {
	if a.Timestamp != b.Timestamp {
		return a.Timestamp > b.Timestamp
	}
	return a.ActorID > b.ActorID 
}