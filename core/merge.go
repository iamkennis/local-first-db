package core

// Merge applies deterministic conflict resolution between two operations.
// Last-write-wins with ID tie-breaking for same timestamps.
func Merge(a, b Operation) Operation {
	if a.Timestamp == 0 {
		return b
	}

	if b.Timestamp > a.Timestamp {
		return b
	}

	if b.Timestamp < a.Timestamp {
		return a
	}

	// Same timestamp → deterministic tie-breaker using ID
	if b.ID > a.ID {
		return b
	}

	return a
}
