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
)

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

	h := handlers.New(st)

	router.Route("/v1", func(r chi.Router) {
		router.Get("/resources", h.ListResources)
		router.Post("/resources", h.CreateResource)
		router.Get("/resources/{id}", h.GetResource)
		//r.Put("/resources/{id}/tags", handler.UpdateResourceTags)
		router.Delete("/resources/{id}", h.DeleteResource)

		router.Post("/resources/{id}/apply-tags", h.ApplyTagsToAzure) //endpoint
	})

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
