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
	pool    map[string]struct{}
	mutex   sync.Mutex
	storage storage.Storage
	jobs    chan *models.Task
}

const maxWorkers = 10

func NewWorker(storage storage.Storage) Worker {
	w := &worker{
		pool:    make(map[string]struct{}),
		storage: storage,
		jobs:    make(chan *models.Task, maxWorkers),
	}

	for i := range maxWorkers {
		go w.newWorker(i, w.jobs)
	}

	return w
}

func (w *worker) Queue(t *models.Task) error {
	w.mutex.Lock()

	if _, ok := w.pool[t.Id]; ok {
		w.mutex.Unlock()
		return errors.New("Task with same id already exists")
	}

	w.pool[t.Id] = struct{}{}
	w.mutex.Unlock()
	w.jobs <- t

	return nil
}

func (w *worker) deleteJob(id string) {
	w.mutex.Lock()
	delete(w.pool, id)
	w.mutex.Unlock()
}

func (w *worker) newWorker(id int, jobs <-chan *models.Task) {
	fmt.Printf("Worker %d waiting for work\n", id)

	for t := range jobs {
		fmt.Printf("Worker %d got new job\n", id)

		err := w.storage.UpdateStatus(t.Id, statuses.InWork)
		if err != nil {
			w.deleteJob(t.Id)
			fmt.Printf("Could not update status = %s of task id = %s", t.Status, t.Id)
			return
		}

		time.Sleep(1 * time.Second)
		fmt.Printf("Task with id = %s finished! Message = %s\n", t.Id, t.Value)

		err = w.storage.UpdateStatus(t.Id, statuses.Done)
		w.deleteJob(t.Id)

		if err != nil {
			fmt.Printf("Could not update status = %s of task id = %s", t.Status, t.Id)
			return
		}
	}
}
