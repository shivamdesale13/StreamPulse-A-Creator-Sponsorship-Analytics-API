package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/streampulse/api/internal/models"
)

type CreatorRepository struct {
	mu   sync.RWMutex
	data map[string]*models.Creator
}

func NewCreatorRepository() *CreatorRepository {
	return &CreatorRepository{data: make(map[string]*models.Creator)}
}

func (r *CreatorRepository) Create(ctx context.Context, creator *models.Creator) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[creator.ID] = creator
	return nil
}

func (r *CreatorRepository) GetByID(ctx context.Context, id string) (*models.Creator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.data[id]
	if !ok {
		return nil, fmt.Errorf("creator %s not found", id)
	}
	return c, nil
}

func (r *CreatorRepository) List(ctx context.Context) ([]*models.Creator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*models.Creator, 0, len(r.data))
	for _, c := range r.data {
		result = append(result, c)
	}
	return result, nil
}

func (r *CreatorRepository) Update(ctx context.Context, creator *models.Creator) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[creator.ID]; !ok {
		return fmt.Errorf("creator %s not found", creator.ID)
	}
	r.data[creator.ID] = creator
	return nil
}

func (r *CreatorRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return fmt.Errorf("creator %s not found", id)
	}
	delete(r.data, id)
	return nil
}
