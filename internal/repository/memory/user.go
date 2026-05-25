package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/streampulse/api/internal/models"
)

type UserRepository struct {
	mu      sync.RWMutex
	byID    map[string]*models.User
	byEmail map[string]*models.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		byID:    make(map[string]*models.User),
		byEmail: make(map[string]*models.User),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byEmail[user.Email]; ok {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}
	r.byID[user.ID] = user
	r.byEmail[user.Email] = user
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byEmail[email]
	if !ok {
		return nil, fmt.Errorf("user with email %s not found", email)
	}
	return u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byID[id]
	if !ok {
		return nil, fmt.Errorf("user %s not found", id)
	}
	return u, nil
}
