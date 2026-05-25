package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/streampulse/api/internal/models"
	"github.com/streampulse/api/internal/queue"
	"github.com/streampulse/api/internal/repository"
)

type DealService struct {
	repo repository.DealRepository
	q    queue.Queue
}

func NewDealService(repo repository.DealRepository, q queue.Queue) *DealService {
	return &DealService{repo: repo, q: q}
}

func (s *DealService) Create(ctx context.Context, req *models.CreateDealRequest) (*models.Deal, error) {
	deal := &models.Deal{
		ID:          uuid.New().String(),
		CreatorID:   req.CreatorID,
		SponsorID:   req.SponsorID,
		Status:      models.DealStatusPending,
		Value:       req.Value,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.repo.Create(ctx, deal); err != nil {
		return nil, err
	}
	_ = s.q.Publish(ctx, queue.DealEvent{
		Type:      queue.EventDealCreated,
		DealID:    deal.ID,
		CreatorID: deal.CreatorID,
		SponsorID: deal.SponsorID,
		Value:     deal.Value,
		Timestamp: time.Now().Format(time.RFC3339),
	})
	return deal, nil
}

func (s *DealService) GetByID(ctx context.Context, id string) (*models.Deal, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DealService) List(ctx context.Context) ([]*models.Deal, error) {
	return s.repo.List(ctx)
}

func (s *DealService) UpdateStatus(ctx context.Context, id string, req *models.UpdateDealStatusRequest) (*models.Deal, error) {
	deal, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !validTransition(deal.Status, req.Status) {
		return nil, fmt.Errorf("invalid transition from %s to %s", deal.Status, req.Status)
	}

	deal.Status = req.Status
	deal.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, deal); err != nil {
		return nil, err
	}

	_ = s.q.Publish(ctx, queue.DealEvent{
		Type:      statusToEvent(req.Status),
		DealID:    deal.ID,
		CreatorID: deal.CreatorID,
		SponsorID: deal.SponsorID,
		Value:     deal.Value,
		Timestamp: time.Now().Format(time.RFC3339),
	})

	return deal, nil
}

func (s *DealService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

var transitions = map[models.DealStatus][]models.DealStatus{
	models.DealStatusPending:   {models.DealStatusActive, models.DealStatusCancelled},
	models.DealStatusActive:    {models.DealStatusCompleted, models.DealStatusCancelled},
	models.DealStatusCompleted: {},
	models.DealStatusCancelled: {},
}

func validTransition(from, to models.DealStatus) bool {
	for _, allowed := range transitions[from] {
		if allowed == to {
			return true
		}
	}
	return false
}

func statusToEvent(status models.DealStatus) queue.EventType {
	switch status {
	case models.DealStatusActive:
		return queue.EventDealActivated
	case models.DealStatusCompleted:
		return queue.EventDealCompleted
	case models.DealStatusCancelled:
		return queue.EventDealCancelled
	default:
		return queue.EventDealUpdated
	}
}
