package internal

import (
	"sync"
)

type Status string

const (
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
)

type Record struct {
	Status     Status
	StatusCode int
	Body       []byte
}

type Store struct {
	mu   sync.Mutex
	data map[string]*Record
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]*Record),
	}
}

func (s *Store) Get(key string) (*Record, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec, ok := s.data[key]
	return rec, ok
}

func (s *Store) SetProcessing(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = &Record{
		Status: StatusProcessing,
	}
}

func (s *Store) SetCompleted(key string, statusCode int, body []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = &Record{
		Status:     StatusCompleted,
		StatusCode: statusCode,
		Body:       body,
	}
}