package memory

import (
	"context"
	"sync"

	"github.com/streampulse/api/internal/models"
)

type AnalyticsRepository struct {
	mu   sync.RWMutex
	data map[string][]*models.Analytics // dealID → entries
}

func NewAnalyticsRepository() *AnalyticsRepository {
	return &AnalyticsRepository{data: make(map[string][]*models.Analytics)}
}

func (r *AnalyticsRepository) Create(ctx context.Context, a *models.Analytics) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[a.DealID] = append(r.data[a.DealID], a)
	return nil
}

func (r *AnalyticsRepository) GetByDealID(ctx context.Context, dealID string) ([]*models.Analytics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.data[dealID], nil
}

func (r *AnalyticsRepository) GetSummary(ctx context.Context) (*models.AnalyticsSummary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	summary := &models.AnalyticsSummary{}
	var convTotal float64
	var entryCount int

	for _, entries := range r.data {
		for _, a := range entries {
			summary.TotalViews += a.Views
			summary.TotalClicks += a.Clicks
			summary.TotalRevenue += a.Revenue
			convTotal += a.ConversionRate
			entryCount++
		}
	}

	if entryCount > 0 {
		summary.AvgConversion = convTotal / float64(entryCount)
	}

	return summary, nil
}
