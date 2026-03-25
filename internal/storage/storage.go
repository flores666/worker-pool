package storage

import (
	"sync"
)

type Storage interface {
	Get(key string) *Task
	Create(t *Task) error
	UpdateStatus(id, status string) error
}

type storage struct {
	m        map[string]*Task
	mutexMap map[string]*sync.Mutex
	mutex    sync.Mutex
}

func NewStorage() Storage {
	return &storage{
		m:        make(map[string]*Task),
		mutexMap: make(map[string]*sync.Mutex),
		mutex:    sync.Mutex{},
	}
}

func (s *storage) Get(key string) *Task {
	m := s.getOrCreateMutex(key)
	m.Lock()
	defer m.Unlock()

	return s.m[key]
}

func (s *storage) Create(t *Task) error {
	m := s.getOrCreateMutex(t.Id)
	m.Lock()
	defer m.Unlock()

	s.m[t.Id] = t

	return nil
}

func (s *storage) UpdateStatus(id, status string) error {
	m := s.getOrCreateMutex(id)
	m.Lock()
	defer m.Unlock()

	s.m[id].Status = status

	return nil
}

func (s *storage) getOrCreateMutex(key string) *sync.Mutex {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if m, ok := s.mutexMap[key]; ok {
		return m
	}

	m := &sync.Mutex{}
	s.mutexMap[key] = m

	return m
}
