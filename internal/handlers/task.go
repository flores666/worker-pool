package handlers

import (
	"net/http"
	"worker-pool/internal/models"
	"worker-pool/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type TasksController struct {
	service service.TasksService
}

func NewController(s service.TasksService) *TasksController {
	return &TasksController{
		service: s,
	}
}

func (c *TasksController) MapRoutes(router chi.Router) {
	router.Post("/task", c.queueTask)
	router.Get("/task/{id}", c.getTask)
}

func (s *TasksController) queueTask(w http.ResponseWriter, r *http.Request) {
	task := &models.Task{}
	err := render.DecodeJSON(r.Body, task)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		return
	}

	s.service.Queue(task)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, task.Id)
}

func (s *TasksController) getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		return
	}

	task := s.service.Get(id)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, task)
}
