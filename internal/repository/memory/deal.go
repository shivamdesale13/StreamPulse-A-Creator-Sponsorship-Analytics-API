package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/streampulse/api/internal/models"
)

type DealRepository struct {
	mu   sync.RWMutex
	data map[string]*models.Deal
}

func NewDealRepository() *DealRepository {
	return &DealRepository{data: make(map[string]*models.Deal)}
}

func (r *DealRepository) Create(ctx context.Context, deal *models.Deal) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[deal.ID] = deal
	return nil
}

func (r *DealRepository) GetByID(ctx context.Context, id string) (*models.Deal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.data[id]
	if !ok {
		return nil, fmt.Errorf("deal %s not found", id)
	}
	return d, nil
}

func (r *DealRepository) List(ctx context.Context) ([]*models.Deal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*models.Deal, 0, len(r.data))
	for _, d := range r.data {
		result = append(result, d)
	}
	return result, nil
}

func (r *DealRepository) ListByCreator(ctx context.Context, creatorID string) ([]*models.Deal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*models.Deal
	for _, d := range r.data {
		if d.CreatorID == creatorID {
			result = append(result, d)
		}
	}
	return result, nil
}

func (r *DealRepository) ListByStatus(ctx context.Context, status models.DealStatus) ([]*models.Deal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*models.Deal
	for _, d := range r.data {
		if d.Status == status {
			result = append(result, d)
		}
	}
	return result, nil
}

func (r *DealRepository) Update(ctx context.Context, deal *models.Deal) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[deal.ID]; !ok {
		return fmt.Errorf("deal %s not found", deal.ID)
	}
	r.data[deal.ID] = deal
	return nil
}

func (r *DealRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return fmt.Errorf("deal %s not found", id)
	}
	delete(r.data, id)
	return nil
}
