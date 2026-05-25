package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/streampulse/api/internal/models"
	"github.com/streampulse/api/internal/repository"
)

type CreatorService struct {
	repo repository.CreatorRepository
}

func NewCreatorService(repo repository.CreatorRepository) *CreatorService {
	return &CreatorService{repo: repo}
}

func (s *CreatorService) Create(ctx context.Context, req *models.CreateCreatorRequest) (*models.Creator, error) {
	creator := &models.Creator{
		ID:              uuid.New().String(),
		Name:            req.Name,
		Platform:        req.Platform,
		ChannelURL:      req.ChannelURL,
		SubscriberCount: req.SubscriberCount,
		Category:        req.Category,
		Email:           req.Email,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := s.repo.Create(ctx, creator); err != nil {
		return nil, err
	}
	return creator, nil
}

func (s *CreatorService) GetByID(ctx context.Context, id string) (*models.Creator, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CreatorService) List(ctx context.Context) ([]*models.Creator, error) {
	return s.repo.List(ctx)
}

func (s *CreatorService) Update(ctx context.Context, id string, req *models.UpdateCreatorRequest) (*models.Creator, error) {
	creator, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		creator.Name = req.Name
	}
	if req.ChannelURL != "" {
		creator.ChannelURL = req.ChannelURL
	}
	if req.SubscriberCount > 0 {
		creator.SubscriberCount = req.SubscriberCount
	}
	if req.Category != "" {
		creator.Category = req.Category
	}
	creator.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, creator); err != nil {
		return nil, err
	}
	return creator, nil
}

func (s *CreatorService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
