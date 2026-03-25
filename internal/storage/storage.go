package storage

import (
	"sync"
)

type Storage interface {
	GetAll() []*Task
	Get(key string) *Task
	Create(t *Task) error
	UpdateStatus(id, status string) error
}

type storage struct {
	m     map[string]*Task
	mutex sync.Mutex
}

func NewStorage() Storage {
	return &storage{
		m: make(map[string]*Task),
	}
}

func (s *storage) Get(key string) *Task {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	t := s.m[key]
	cp := *t

	return &cp
}

func (s *storage) GetAll() []*Task {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	result := make([]*Task, 0, len(s.m))
	for _, t := range s.m {
		cp := *t
		result = append(result, &cp)
	}

	return result
}

func (s *storage) Create(t *Task) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.m[t.Id] = t

	return nil
}

func (s *storage) UpdateStatus(id, status string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.m[id].Status = status

	return nil
}
