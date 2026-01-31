package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ThiagoScheffer/azure-tagger-api/internal/store"
	"github.com/go-chi/chi/v5"
)

type mockTagger struct {
	called     bool
	resourceID string
	tags       map[string]string
	err        error
}

func (m *mockTagger) ApplyTags(ctx context.Context, resourceID string, tags map[string]string) error {
	m.called = true
	m.resourceID = resourceID
	m.tags = tags
	return m.err
}

func newTestRouterWithApply(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Post("/resources", h.CreateResource)
		r.Post("/resources/{id}/apply-tags", h.ApplyTagsToAzure)
	})
	return r
}

func TestHandlers_ApplyTagsToAzure_Success(t *testing.T) {
	st := store.NewMemoryStore()
	h := New(st)

	mt := &mockTagger{}
	h.taggerFactory = func() (AzureTagger, error) { return mt, nil }

	router := newTestRouterWithApply(h)

	// Create a resource first
	created := st.Create(
		"vm-1",
		"/subscriptions/x/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/vm-1",
		map[string]string{},
	)

	body := map[string]any{
		"tags": map[string]string{"owner": "jairo", "project": "portfolio"},
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/v1/resources/"+created.ID+"/apply-tags", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !mt.called {
		t.Fatal("expected mock tagger to be called")
	}
	if mt.resourceID != created.AzureID {
		t.Fatalf("expected resourceID %q, got %q", created.AzureID, mt.resourceID)
	}
	if mt.tags["owner"] != "jairo" {
		t.Fatalf("expected tag owner=jairo, got %v", mt.tags["owner"])
	}
}

func TestHandlers_ApplyTagsToAzure_Validation(t *testing.T) {
	st := store.NewMemoryStore()
	h := New(st)

	mt := &mockTagger{}
	h.taggerFactory = func() (AzureTagger, error) { return mt, nil }

	router := newTestRouterWithApply(h)

	created := st.Create("vm-1", "/subscriptions/x/.../vm-1", map[string]string{})

	req := httptest.NewRequest(http.MethodPost, "/v1/resources/"+created.ID+"/apply-tags", bytes.NewBufferString(`{"tags":{}}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if mt.called {
		t.Fatal("did not expect mock tagger to be called on validation failure")
	}
}
