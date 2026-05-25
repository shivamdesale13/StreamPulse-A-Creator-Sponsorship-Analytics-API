package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/streampulse/api/internal/models"
	"github.com/streampulse/api/internal/repository"
)

type AnalyticsService struct {
	analyticsRepo repository.AnalyticsRepository
	dealRepo      repository.DealRepository
}

func NewAnalyticsService(analyticsRepo repository.AnalyticsRepository, dealRepo repository.DealRepository) *AnalyticsService {
	return &AnalyticsService{analyticsRepo: analyticsRepo, dealRepo: dealRepo}
}

func (s *AnalyticsService) RecordForDeal(ctx context.Context, dealID string, req *models.CreateAnalyticsRequest) (*models.Analytics, error) {
	if _, err := s.dealRepo.GetByID(ctx, dealID); err != nil {
		return nil, err
	}

	var convRate float64
	if req.Views > 0 {
		convRate = float64(req.Clicks) / float64(req.Views) * 100
	}

	a := &models.Analytics{
		ID:             uuid.New().String(),
		DealID:         dealID,
		Views:          req.Views,
		Clicks:         req.Clicks,
		ConversionRate: convRate,
		Revenue:        req.Revenue,
		RecordedAt:     time.Now(),
	}

	if err := s.analyticsRepo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *AnalyticsService) GetForDeal(ctx context.Context, dealID string) ([]*models.Analytics, error) {
	return s.analyticsRepo.GetByDealID(ctx, dealID)
}

func (s *AnalyticsService) GetCreatorAnalytics(ctx context.Context, creatorID string) ([]*models.Analytics, error) {
	deals, err := s.dealRepo.ListByCreator(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	var all []*models.Analytics
	for _, deal := range deals {
		entries, err := s.analyticsRepo.GetByDealID(ctx, deal.ID)
		if err != nil {
			continue
		}
		all = append(all, entries...)
	}
	return all, nil
}

func (s *AnalyticsService) GetSummary(ctx context.Context) (*models.AnalyticsSummary, error) {
	summary, err := s.analyticsRepo.GetSummary(ctx)
	if err != nil {
		return nil, err
	}

	active, _ := s.dealRepo.ListByStatus(ctx, models.DealStatusActive)
	completed, _ := s.dealRepo.ListByStatus(ctx, models.DealStatusCompleted)
	all, _ := s.dealRepo.List(ctx)

	summary.ActiveDeals = len(active)
	summary.CompletedDeals = len(completed)
	summary.TotalDeals = len(all)

	return summary, nil
}
