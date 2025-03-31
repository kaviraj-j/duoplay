package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/kaviraj-j/duoplay/internal/model"
)

var (
	ErrUserNotFound error = fmt.Errorf("user not found")
)

// UserRepository defines the methods for User related operations
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
}

// inMemoryUserRepository implements UserRepository interface and stores the user data within the app memory
type inMemoryUserRepository struct {
	users map[string]*model.User
	mu    sync.RWMutex
}

// NewUserRepository creates a new in-memory user repository
func NewUserRepository() UserRepository {
	return &inMemoryUserRepository{
		users: make(map[string]*model.User),
	}
}

// implement UserRepository methods
func (repository *inMemoryUserRepository) Create(ctx context.Context, user *model.User) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	repository.users[user.ID] = user
	return nil
}

func (repository *inMemoryUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	user, ok := repository.users[id]
	if !ok {
		return &model.User{}, ErrUserNotFound
	}
	return user, nil
}
