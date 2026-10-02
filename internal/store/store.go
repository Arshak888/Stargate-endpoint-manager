package store

import (
	"crypto/subtle"
	"errors"
	"sync"
	"time"

	"github.com/Arshak888/Stargate-endpoint-manager/internal/auth"
	"github.com/Arshak888/Stargate-endpoint-manager/internal/domain"
)

var ErrNotFound = errors.New("not found")

type MemoryStore struct {
	mu           sync.RWMutex
	enrollments  map[string]domain.Enrollment
	endpoints    map[string]domain.Endpoint
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		enrollments: map[string]domain.Enrollment{},
		endpoints:   map[string]domain.Endpoint{},
	}
}

func (s *MemoryStore) PutEnrollment(e domain.Enrollment) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.enrollments[e.ID] = e
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
	return e, nil
}

func (s *MemoryStore) PutEndpoint(e domain.Endpoint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.endpoints[e.ID] = e
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
