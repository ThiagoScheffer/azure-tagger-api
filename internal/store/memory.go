package store

import (
	"errors"
	"sync"
	"time"

	"github.com/ThiagoScheffer/azure-tagger-api/internal/models"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("not found !")

type MemoryStore struct {
	mu        sync.RWMutex
	resources map[string]models.Resource
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		resources: make(map[string]models.Resource), // init the map and return the struct!
	}
}

// Create a new resource in the store and return it !!
func (s *MemoryStore) Create(name, azureID string, tags map[string]string) models.Resource {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.NewString()
	r := models.Resource{
		ID:          id,
		Name:        name,
		AzureID:     azureID,
		Tags:        tags,
		CreatedUnix: time.Now().Unix(),
	}
	s.resources[id] = r
	return r
}

func (s *MemoryStore) List() []models.Resource {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]models.Resource, 0, len(s.resources))
	for _, v := range s.resources {
		out = append(out, v)
	}
	return out
}

func (s *MemoryStore) Get(id string) (models.Resource, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.resources[id]
	if !ok {
		return models.Resource{}, ErrNotFound
	}
	return v, nil
}

func (s *MemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.resources[id]; !ok {
		return ErrNotFound
	}
	delete(s.resources, id)
	return nil
}
