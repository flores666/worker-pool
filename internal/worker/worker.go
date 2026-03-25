package worker

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"worker-pool/internal/models"
	"worker-pool/internal/statuses"
	"worker-pool/internal/storage"
)

type Worker interface {
	Queue(t *models.Task) error
}

type worker struct {
	pool    map[string]string
	mutex   sync.Mutex
	storage storage.Storage
}

func NewWorker(storage storage.Storage) Worker {
	return &worker{
		pool:    make(map[string]string),
		storage: storage,
	}
}

func (w *worker) Queue(t *models.Task) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if _, ok := w.pool[t.Id]; ok {
		return errors.New("Task with same id already exists")
	}

	w.pool[t.Id] = t.Status
	go w.doWork(t)

	return nil
}

func (w *worker) doWork(t *models.Task) {
	err := w.storage.UpdateStatus(t.Id, statuses.InWork)
	if err != nil {
		fmt.Printf("Could not update status = %s of task id = %s", t.Status, t.Id)
		return
	}

	time.Sleep(8 * time.Second)
	fmt.Printf("Task with id = %s finished! Message = %s\n", t.Id, t.Value)

	err = w.storage.UpdateStatus(t.Id, statuses.Done)
	if err != nil {
		fmt.Printf("Could not update status = %s of task id = %s", t.Status, t.Id)
		return
	}
}
