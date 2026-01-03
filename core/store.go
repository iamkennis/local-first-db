package core

type Store struct {
	ops     []Operation
	seen    map[string]bool
	input   chan Operation
	storage Storage
}

// NewStore creates a new store with the given storage backend
func NewStore(storage Storage) (*Store, error) {
	ops, err := storage.Load()
	if err != nil {
		return nil, err
	}

	s := &Store{
		ops:     ops,
		seen:    map[string]bool{},
		input:   make(chan Operation),
		storage: storage,
	}

	go s.run()
	return s, nil
}

// run processes operations in a single writer pattern
func (s *Store) run() {
	for op := range s.input {
		if s.seen[op.ID] {
			continue
		}
		s.seen[op.ID] = true
		_ = s.storage.Append(op)
		s.ops = append(s.ops, op)
	}
}

// Apply adds an operation to the store
func (s *Store) Apply(op Operation) {
	s.input <- op
}

// Get retrieves the latest value for a key using CRDT merge logic
func (s *Store) Get(key string) []byte {
	var latest Operation
	for _, op := range s.ops {
		if op.Key == key {
			latest = Merge(latest, op)
		}
	}
	if latest.Type == "delete" {
		return nil
	}
	return latest.Value
}