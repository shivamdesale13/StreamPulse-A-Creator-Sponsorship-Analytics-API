package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/streampulse/api/internal/models"
)

type SponsorRepository struct {
	mu   sync.RWMutex
	data map[string]*models.Sponsor
}

func NewSponsorRepository() *SponsorRepository {
	return &SponsorRepository{data: make(map[string]*models.Sponsor)}
}

func (r *SponsorRepository) Create(ctx context.Context, sponsor *models.Sponsor) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[sponsor.ID] = sponsor
	return nil
}

func (r *SponsorRepository) GetByID(ctx context.Context, id string) (*models.Sponsor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.data[id]
	if !ok {
		return nil, fmt.Errorf("sponsor %s not found", id)
	}
	return s, nil
}

func (r *SponsorRepository) List(ctx context.Context) ([]*models.Sponsor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*models.Sponsor, 0, len(r.data))
	for _, s := range r.data {
		result = append(result, s)
	}
	return result, nil
}

func (r *SponsorRepository) Update(ctx context.Context, sponsor *models.Sponsor) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[sponsor.ID]; !ok {
		return fmt.Errorf("sponsor %s not found", sponsor.ID)
	}
	r.data[sponsor.ID] = sponsor
	return nil
}

func (r *SponsorRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return fmt.Errorf("sponsor %s not found", id)
	}
	delete(r.data, id)
	return nil
}
