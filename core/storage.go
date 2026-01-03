package core

// Storage defines the interface for persisting operations
type Storage interface {
	Append(op Operation) error
	Load() ([]Operation, error)
}

