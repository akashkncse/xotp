package otp

import "sync"

type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewStore() *Store {
	return &Store{data: make(map[string]string)}
}

func (s *Store) Save(userID, codeHash string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[userID] = codeHash
}

func (s *Store) Find(userID string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	codeHash, exists := s.data[userID]
	return codeHash, exists
}

func (s *Store) Delete(userID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[userID]; !exists {
		return false
	}

	delete(s.data, userID)
	return true
}
