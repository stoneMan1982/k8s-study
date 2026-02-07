package friend

import "sync"

// Service holds the internal state for the friend module.
type Service struct {
	mu      sync.RWMutex
	started bool
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.started = true
	return nil
}

func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.started = false
	return nil
}

func (s *Service) Started() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.started
}
