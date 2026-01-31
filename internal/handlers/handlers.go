package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/ThiagoScheffer/azure-tagger-api/internal/azure"
	"github.com/ThiagoScheffer/azure-tagger-api/internal/store"
	"github.com/go-chi/chi/v5"
	//"github.com/ThiagoScheffer/azure-tagger-api/internal/models"
)

type Handler struct {
	store *store.MemoryStore

	taggerFactory TaggerFactory //for testing
}

func New(st *store.MemoryStore) *Handler {
	return &Handler{
		store: st,
		taggerFactory: func() (AzureTagger, error) {
			apiVersion := os.Getenv("AZURE_RESOURCE_API_VERSION")
			// example: "2021-04-01" (depends on resource type!)
			return azure.NewTagger(apiVersion)
		},
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

type createReq struct {
	Name    string            `json:"name"`
	AzureID string            `json:"azureId"`
	Tags    map[string]string `json:"tags"`
}

// CreateResource godoc
// @Summary      Create a resource
// @Description  Stores an Azure resource ID + tags
// @Tags         resources
// @Accept       json
// @Produce      json
// @Param        payload  body      createReq  true  "Resource payload"
// @Success      201      {object}  models.Resource
// @Failure      400      {object}  map[string]string
// @Router       /resources [post]
func (h *Handler) CreateResource(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	if req.Name == "" || req.AzureID == "" {
		writeErr(w, 400, "name and azureId are required")
		return
	}
	if req.Tags == nil {
		req.Tags = map[string]string{}
	}
	res := h.store.Create(req.Name, req.AzureID, req.Tags)
	writeJSON(w, 201, res)
}

func (h *Handler) ListResources(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, h.store.List())
}

func (h *Handler) GetResource(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res, err := h.store.Get(id)
	if err != nil {
		writeErr(w, 404, "not found")
		return
	}
	writeJSON(w, 200, res)
}

func (h *Handler) DeleteResource(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.store.Delete(id); err != nil {
		writeErr(w, 404, "not found")
		return
	}
	w.WriteHeader(204)
}

type AzureTagger interface {
	ApplyTags(ctx context.Context, resourceID string, tags map[string]string) error
}

type TaggerFactory func() (AzureTagger, error)
