package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ThiagoScheffer/azure-tagger-api/internal/azure"
	"github.com/go-chi/chi/v5"
)

type applyReq struct {
	Tags map[string]string `json:"tags"`
}

func (h *Handler) ApplyTagsToAzure(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.store.Get(id)
	if err != nil {
		writeErr(w, 404, "not found")
		return
	}

	var req applyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	if len(req.Tags) == 0 {
		writeErr(w, 400, "tags required")
		return
	}

	tagger, err := azure.NewTagger()
	if err != nil {
		writeErr(w, 400, "azure not configured: set AZURE_SUBSCRIPTION_ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := tagger.ApplyTags(ctx, res.AzureID, "2020-01-01", req.Tags); err != nil {
		writeErr(w, 500, "azure error: "+err.Error())
		return
	}

	writeJSON(w, 200, map[string]any{
		"message":  "tags applied",
		"resource": res.AzureID,
		"tags":     req.Tags,
	})
}
