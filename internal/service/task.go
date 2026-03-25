package service

import (
	"errors"
	"worker-pool/internal/models"
	"worker-pool/internal/statuses"
	"worker-pool/internal/storage"
	"worker-pool/internal/worker"

	"github.com/google/uuid"
)

type TasksService interface {
	GetAll() []*models.Task
	Get(key string) *models.Task
	Queue(t *models.Task) (string, error)
}

type service struct {
	storage   storage.Storage
	validator Validator
	worker    worker.Worker
}

func NewTasksService(s storage.Storage, validator Validator, worker worker.Worker) TasksService {
	return &service{
		storage:   s,
		validator: validator,
		worker:    worker,
	}
}

func (s *service) GetAll() []*models.Task {
	all := s.storage.GetAll()
	result := make([]*models.Task, 0, len(all))

	for _, item := range all {
		result = append(result, &models.Task{
			Id:     item.Id,
			Status: item.Status,
			Value:  item.Value,
		})
	}

	return result
}

func (s *service) Get(key string) *models.Task {
	if key == "" {
		return nil
	}

	t := (s.storage.Get(key))
	if t == nil {
		return nil
	}

	return &models.Task{
		Id:     t.Id,
		Status: t.Status,
	}
}

func (s *service) Queue(t *models.Task) (string, error) {
	if t.Id == "" {
		t.Id = uuid.New().String()
	}

	t.Status = statuses.Draft

	err := s.storage.Create(&storage.Task{
		Id:     t.Id,
		Status: t.Status,
		Value:  t.Value,
	})

	if err != nil {
		return "", errors.New("Error while saving task to storage")
	}

	err = s.worker.Queue(t)
	if err != nil {
		return t.Id, errors.New("Error while adding task to queue")
	}

	return t.Id, nil
}
