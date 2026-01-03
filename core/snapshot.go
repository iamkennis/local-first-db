package core

type Snapshot struct {
	LastTimestamp int64
	State         map[string]string
}