package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/streampulse/api/internal/models"
	"github.com/streampulse/api/internal/repository"
)

type SponsorService struct {
	repo repository.SponsorRepository
}

func NewSponsorService(repo repository.SponsorRepository) *SponsorService {
	return &SponsorService{repo: repo}
}

func (s *SponsorService) Create(ctx context.Context, req *models.CreateSponsorRequest) (*models.Sponsor, error) {
	sponsor := &models.Sponsor{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Industry:     req.Industry,
		Website:      req.Website,
		ContactEmail: req.ContactEmail,
		Budget:       req.Budget,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := s.repo.Create(ctx, sponsor); err != nil {
		return nil, err
	}
	return sponsor, nil
}

func (s *SponsorService) GetByID(ctx context.Context, id string) (*models.Sponsor, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SponsorService) List(ctx context.Context) ([]*models.Sponsor, error) {
	return s.repo.List(ctx)
}

func (s *SponsorService) Update(ctx context.Context, id string, req *models.UpdateSponsorRequest) (*models.Sponsor, error) {
	sponsor, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		sponsor.Name = req.Name
	}
	if req.Website != "" {
		sponsor.Website = req.Website
	}
	if req.Budget > 0 {
		sponsor.Budget = req.Budget
	}
	sponsor.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, sponsor); err != nil {
		return nil, err
	}
	return sponsor, nil
}

func (s *SponsorService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
