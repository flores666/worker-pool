package controllers

import (
	"net/http"
	"worker-pool/internal/storage"
	"worker-pool/internal/tasks"
	"worker-pool/internal/validator"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type TasksController struct {
	storage   storage.Storage
	validator validator.Validator
}

func NewController(storage storage.Storage, validator validator.Validator) *TasksController {
	return &TasksController{
		storage:   storage,
		validator: validator,
	}
}

func (c *TasksController) MapRoutes(router chi.Router) {
	router.Post("/task", c.createTask)
	router.Post("/task/{id}", c.getTask)
}

func (s *TasksController) createTask(w http.ResponseWriter, r *http.Request) {
	task := &tasks.Task{}
	err := render.DecodeJSON(r.Body, task)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		return
	}

	s.storage.Create(task)
	render.Status(r, http.StatusOK)
}

func (s *TasksController) getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		return
	}

	task := s.storage.Get(id)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, task)
}
