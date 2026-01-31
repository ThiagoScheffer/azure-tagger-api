package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/ThiagoScheffer/azure-tagger-api/internal/store"
)

func newTestRouter(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Post("/resources", h.CreateResource)
		r.Get("/resources", h.ListResources)
		r.Get("/resources/{id}", h.GetResource)
		r.Delete("/resources/{id}", h.DeleteResource)
	})
	return r
}

func TestHandlers_Create_List_Get_Delete_HappyPath(t *testing.T) {
	st := store.NewMemoryStore()
	h := New(st)
	router := newTestRouter(h)

	// Create
	createBody := map[string]any{
		"name":    "vm-1",
		"azureId": "/subscriptions/x/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/vm-1",
		"tags": map[string]string{
			"env": "dev",
		},
	}
	b, _ := json.Marshal(createBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/resources", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", rr.Code, rr.Body.String())
	}

	var created map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &created); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	id, ok := created["id"].(string)
	if !ok || id == "" {
		t.Fatalf("expected id in response, got: %v", created["id"])
	}

	// List
	req = httptest.NewRequest(http.MethodGet, "/v1/resources", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var list []map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &list); err != nil {
		t.Fatalf("invalid json list: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(list))
	}

	// Get
	req = httptest.NewRequest(http.MethodGet, "/v1/resources/"+id, nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	// Delete
	req = httptest.NewRequest(http.MethodDelete, "/v1/resources/"+id, nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}

	// Get after delete -> 404
	req = httptest.NewRequest(http.MethodGet, "/v1/resources/"+id, nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestHandlers_Create_TableDriven(t *testing.T) {
	st := store.NewMemoryStore()
	h := New(st)
	router := newTestRouter(h)

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "invalid json",
			body:       "{invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing name",
			body:       `{"azureId":"/subscriptions/x/...","tags":{"a":"b"}}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing azureId",
			body:       `{"name":"vm-1","tags":{"a":"b"}}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "ok",
			body:       `{"name":"vm-1","azureId":"/subscriptions/x/.../vm-1","tags":{"env":"dev"}}`,
			wantStatus: http.StatusCreated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/v1/resources", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("expected %d, got %d, body=%s", tc.wantStatus, rr.Code, rr.Body.String())
			}
		})
	}
}
