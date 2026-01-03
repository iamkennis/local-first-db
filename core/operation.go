package core

type Operation struct {
	ID        string `json:"id"`
	Actor     string `json:"actor"`
	Timestamp int64  `json:"ts"`
	Type      string `json:"type"` // set | delete
	Key       string `json:"key"`
	Value     []byte `json:"value"`
}