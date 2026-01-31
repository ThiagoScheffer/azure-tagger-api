package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type applyReq struct {
	Tags map[string]string `json:"tags"`
}

// ApplyTagsToAzure godoc
// @Summary      Apply tags to the Azure resource
// @Tags         azure
// @Accept       json
// @Produce      json
// @Param        id      path     string   true  "Resource ID"
// @Param        payload body     applyReq true  "Tags to apply"
// @Success      200     {object} map[string]any
// @Failure      400     {object} map[string]string
// @Failure      404     {object} map[string]string
// @Failure      500     {object} map[string]string
// @Router       /resources/{id}/apply-tags [post]
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

	tagger, err := h.taggerFactory() //for testing
	if err != nil {
		writeErr(w, 400, "azure not configured: set AZURE_SUBSCRIPTION_ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := tagger.ApplyTags(ctx, res.AzureID, req.Tags); err != nil {
		writeErr(w, 500, "azure error: "+err.Error())
		return
	}

	writeJSON(w, 200, map[string]any{
		"message":  "tags applied",
		"resource": res.AzureID,
		"tags":     req.Tags,
	})
}
