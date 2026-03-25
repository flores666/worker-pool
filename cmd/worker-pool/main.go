package main

import (
	"fmt"
	"net/http"
	"worker-pool/internal/handlers"
	"worker-pool/internal/service"
	"worker-pool/internal/storage"
	"worker-pool/internal/worker"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	router := chi.NewRouter()
	fmt.Println("Router Created")

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	storage := storage.NewStorage()
	tasksController := handlers.NewController(service.NewTasksService(storage, service.NewValidator(), worker.NewWorker(storage)))
	fmt.Println("Api Handlers Registered")

	tasksController.MapRoutes(router)
	fmt.Println("Routes Mapped")

	fmt.Println("Server Started! Address: http://localhost:5124")
	err := http.ListenAndServe(":5124", router)
	if err != nil {
		fmt.Printf("ERROR = %s\n", err.Error())
	}
}
