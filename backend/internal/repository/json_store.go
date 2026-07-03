package repository

import (
	"encoding/json"
	"os"
	"sync"
)

type JSONStore struct {
	mu sync.RWMutex
}

func NewJSONStore() *JSONStore {
	return &JSONStore{}
}

func (s *JSONStore) Read(path string, v any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func (s *JSONStore) Write(path string, v any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
