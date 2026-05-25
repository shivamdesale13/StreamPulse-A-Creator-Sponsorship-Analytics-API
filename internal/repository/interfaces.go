package repository

import (
	"context"

	"github.com/streampulse/api/internal/models"
)

type CreatorRepository interface {
	Create(ctx context.Context, creator *models.Creator) error
	GetByID(ctx context.Context, id string) (*models.Creator, error)
	List(ctx context.Context) ([]*models.Creator, error)
	Update(ctx context.Context, creator *models.Creator) error
	Delete(ctx context.Context, id string) error
}

type SponsorRepository interface {
	Create(ctx context.Context, sponsor *models.Sponsor) error
	GetByID(ctx context.Context, id string) (*models.Sponsor, error)
	List(ctx context.Context) ([]*models.Sponsor, error)
	Update(ctx context.Context, sponsor *models.Sponsor) error
	Delete(ctx context.Context, id string) error
}

type DealRepository interface {
	Create(ctx context.Context, deal *models.Deal) error
	GetByID(ctx context.Context, id string) (*models.Deal, error)
	List(ctx context.Context) ([]*models.Deal, error)
	ListByCreator(ctx context.Context, creatorID string) ([]*models.Deal, error)
	ListByStatus(ctx context.Context, status models.DealStatus) ([]*models.Deal, error)
	Update(ctx context.Context, deal *models.Deal) error
	Delete(ctx context.Context, id string) error
}

type AnalyticsRepository interface {
	Create(ctx context.Context, analytics *models.Analytics) error
	GetByDealID(ctx context.Context, dealID string) ([]*models.Analytics, error)
	GetSummary(ctx context.Context) (*models.AnalyticsSummary, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
}
