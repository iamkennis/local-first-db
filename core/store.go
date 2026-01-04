package core

import "time"

type Store struct {
	data map[string]Operation
	log  []Operation
	snapshot *Snapshot
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]Operation),
	}
}


func (s *Store) Apply(op Operation) {
	if s.snapshot != nil && op.Timestamp <= s.snapshot.LastTimestamp {
		return 
	}

	existing, ok := s.data[op.Key]
	if !ok || newer(op, existing) {
	s.data[op.Key] = op
}
	s.log = append(s.log, op)
}

func (s *Store) State() map[string]string {
	out := map[string]string{}
	for k, v := range s.data {
		out[k] = v.Value
	}
	return out
}

func (s *Store) Ops() []Operation {
	return s.log
}

func (s *Store) CreateSnapshot() Snapshot {
	state := map[string]string{}
	for k, v := range s.data {
		state[k] = v.Value
	}

	snap := Snapshot{
		LastTimestamp: time.Now().UnixNano(),
		State:         state,
	}
	s.snapshot = &snap
	return snap
}