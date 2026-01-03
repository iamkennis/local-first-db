package core

import (
	"crypto/rand"

	"github.com/google/uuid"
)

type Identity struct {
	DeviceID string
	Key      []byte
}

func NewIdentity() (*Identity, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return &Identity{
		DeviceID: uuid.NewString(),
		Key:      key,
	}, nil
}