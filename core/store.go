package core

type Store struct {
	data map[string]Operation
	log  []Operation
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]Operation),
	}
}

func (s *Store) Apply(op Operation) {
	existing, ok := s.data[op.Key]
	if !ok || op.Timestamp > existing.Timestamp {
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