package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ThiagoScheffer/azure-tagger-api/internal/handlers"
	"github.com/ThiagoScheffer/azure-tagger-api/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/ThiagoScheffer/azure-tagger-api/docs"
)

// infos for swaggrs
// @title           Azure Tagger API
// @version         1.0
// @description     REST API to register Azure resource IDs and apply tags to them.
// @termsOfService  https://example.com/terms
// @contact.name    Your Name
// @contact.email   you@example.com
// @license.name    MIT
// @host            localhost:8080
// @BasePath        /v1

// @description REST API that manages Azure resource metadata and applies tags using Azure SDK.
// @description This project demonstrates:
// @description - Clean architecture
// @description - Dependency injection
// @description - Unit testing with mocks
// @description - Docker containerization
// @description - CI/CD with GitHub Actions

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	st := store.NewMemoryStore()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	router.Get("/swagger/*", httpSwagger.WrapHandler) //for swagger ui

	h := handlers.New(st)

	router.Route("/v1", func(r chi.Router) {

		r.Post("/resources", h.CreateResource)

		// @Summary Get a resource
		// @Tags    resources
		// @Produce json
		// @Param   id   path     string true "Resource ID"
		// @Success 200  {object} models.Resource
		// @Failure 404  {object} map[string]string
		// @Router  /resources/{id} [get]
		r.Get("/resources", h.ListResources)

		// @Summary Get a resource
		// @Tags    resources
		// @Produce json
		// @Param   id   path     string true "Resource ID"
		// @Success 200  {object} models.Resource
		// @Failure 404  {object} map[string]string
		// @Router  /resources/{id} [get]
		r.Get("/resources/{id}", h.GetResource)
		//r.Put("/resources/{id}/tags", handler.UpdateResourceTags)

		// @Summary Delete a resource
		// @Tags    resources
		// @Param   id path string true "Resource ID"
		// @Success 204
		// @Failure 404 {object} map[string]string
		// @Router  /resources/{id} [delete]
		r.Delete("/resources/{id}", h.DeleteResource)

		r.Post("/resources/{id}/apply-tags", h.ApplyTagsToAzure) //endpoint
	})

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
