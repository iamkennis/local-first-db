package file

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/iamkennis/decentralized-db/core"
)

type FileStorage struct {
	path string
}

func New(path string) *FileStorage {
	return &FileStorage{path}
}

func (f *FileStorage) Append(op core.Operation) error {
	file, _ := os.OpenFile(f.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer file.Close()

	data, _ := json.Marshal(op)
	_, err := file.Write(append(data, '\n'))
	return err
}

func (f *FileStorage) Load() ([]core.Operation, error) {
	file, err := os.Open(f.path)
	if os.IsNotExist(err) {
		return []core.Operation{}, nil
	}

	var ops []core.Operation
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var op core.Operation
		json.Unmarshal(scanner.Bytes(), &op)
		ops = append(ops, op)
	}
	return ops, nil
}