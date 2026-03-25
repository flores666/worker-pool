package storage

import (
	"sync"
	"worker-pool/internal/tasks"

	"github.com/google/uuid"
)

type Storage interface {
	Get(key string) *tasks.Task
	Create(t *tasks.Task)
}

type storage struct {
	m        map[string]*tasks.Task
	mutexMap map[string]*sync.Mutex
	mutex    sync.Mutex
}

func NewStorage() Storage {
	return &storage{
		m:        make(map[string]*tasks.Task),
		mutexMap: make(map[string]*sync.Mutex),
		mutex:    sync.Mutex{},
	}
}

func (s *storage) Get(key string) *tasks.Task {
	/*m := s.getOrCreateMutex(key)
	m.Lock()
	defer m.Unlock()*/

	return s.m[key]
}

func (s *storage) Create(t *tasks.Task) {
	if t.Id == "" {
		key := uuid.New().String()
		t.Id = key
	}

	s.m[t.Id] = t
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
