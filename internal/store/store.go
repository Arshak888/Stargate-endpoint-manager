package store

import "sync"

type MemoryStore struct {
	mu sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}
