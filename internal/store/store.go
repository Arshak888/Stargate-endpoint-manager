package store

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Arshak888/Stargate-endpoint-manager/internal/auth"
	"github.com/Arshak888/Stargate-endpoint-manager/internal/domain"
)

var ErrNotFound = errors.New("not found")

type diskState struct {
	Enrollments map[string]domain.Enrollment `json:"enrollments"`
	Endpoints   map[string]domain.Endpoint   `json:"endpoints"`
}

type MemoryStore struct {
	mu          sync.RWMutex
	path        string
	enrollments map[string]domain.Enrollment
	endpoints   map[string]domain.Endpoint
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		enrollments: map[string]domain.Enrollment{},
		endpoints:   map[string]domain.Endpoint{},
	}
}

func NewPersistentStore(path string) (*MemoryStore, error) {
	s := NewMemoryStore()
	s.path = path
	if path == "" {
		return s, nil
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}

	var disk diskState
	if err := json.Unmarshal(data, &disk); err != nil {
		return nil, err
	}
	if disk.Enrollments != nil {
		s.enrollments = disk.Enrollments
	}
	if disk.Endpoints != nil {
		s.endpoints = disk.Endpoints
	}
	return s, nil
}

func (s *MemoryStore) persistLocked() error {
	if s.path == "" {
		return nil
	}

	data, err := json.MarshalIndent(diskState{
		Enrollments: s.enrollments,
		Endpoints:   s.endpoints,
	}, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func (s *MemoryStore) PutEnrollment(e domain.Enrollment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.enrollments[e.ID] = e
	return s.persistLocked()
}

func (s *MemoryStore) GetEnrollment(id string) (domain.Enrollment, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.enrollments[id]
	return e, ok
}

func (s *MemoryStore) ConsumeEnrollment(id string, now time.Time) (domain.Enrollment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.enrollments[id]
	if !ok || e.UsedAt != nil || !now.Before(e.ExpiresAt) {
		return domain.Enrollment{}, ErrNotFound
	}
	e.UsedAt = &now
	s.enrollments[id] = e
	if err := s.persistLocked(); err != nil {
		return domain.Enrollment{}, err
	}
	return e, nil
}

func (s *MemoryStore) PutEndpoint(e domain.Endpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.endpoints[e.ID] = e
	return s.persistLocked()
}

func (s *MemoryStore) GetEndpoint(id string) (domain.Endpoint, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.endpoints[id]
	return e, ok
}

func (s *MemoryStore) AuthenticateEndpoint(id, credential string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ep, ok := s.endpoints[id]
	if !ok || ep.CredentialHash == "" || credential == "" {
		return false
	}
	return subtle.ConstantTimeCompare(
		[]byte(ep.CredentialHash),
		[]byte(auth.HashToken(credential)),
	) == 1
}
