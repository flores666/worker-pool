package main

import (
	"net/http"
	"worker-pool/internal/controllers"
	"worker-pool/internal/storage"
	"worker-pool/internal/validator"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	tasksController := controllers.NewController(storage.NewStorage(), validator.NewValidator())
	tasksController.MapRoutes(router)

	http.ListenAndServe(":3000", router)
}
