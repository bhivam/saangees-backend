package data

import (
	"errors"
	"sync"
	"time"
)

// InMemorySessionStore implements the SessionStore interface with an in-memory array.
type InMemorySessionStore struct {
	sessions []*Session
	mu       sync.Mutex // Protects access to the sessions slice
}

// NewInMemorySessionStore creates and returns a new instance of InMemorySessionStore.
func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		sessions: make([]*Session, 0),
	}
}

// CreateSession adds a new session to the in-memory store.
func (s *InMemorySessionStore) CreateSession(session *Session) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session.CreatedAt = time.Now()
	s.sessions = append(s.sessions, session)
	return session, nil
}

// GetSession retrieves a session by ID.
func (s *InMemorySessionStore) GetSession(id string) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, session := range s.sessions {
		if session.ID == id {
			return session, nil
		}
	}
	return nil, errors.New("session not found")
}

// RevokedSession marks a session as revoked by ID.
func (s *InMemorySessionStore) RevokedSession(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, session := range s.sessions {
		if session.ID == id {
			session.IsRevoked = true
			return nil
		}
	}
	return errors.New("session not found")
}

// DeleteSession removes a session by ID from the in-memory store.
func (s *InMemorySessionStore) DeleteSession(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, session := range s.sessions {
		if session.ID == id {
			// Remove the session from the slice
			s.sessions = append(s.sessions[:i], s.sessions[i+1:]...)
			return nil
		}
	}
	return errors.New("session not found")
}
