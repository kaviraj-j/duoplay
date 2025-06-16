package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/kaviraj-j/duoplay/internal/model"
)

type GameRepository interface {
	GetGamesList(ctx context.Context) ([]string, error)
}

type inMemoryGameRepository struct {
	games map[string]struct {
		Game *model.Game
		Name string
	}
	mu sync.RWMutex
}

var (
	ErrGameNotFound error = fmt.Errorf("game not found")
)

func NewGameRepository() GameRepository {
	return &inMemoryGameRepository{
		games: make(map[string]struct {
			Game *model.Game
			Name string
		}),
	}
}

func (gameRepository *inMemoryGameRepository) GetGamesList(ctx context.Context) ([]string, error) {
	gameRepository.mu.RLock()
	defer gameRepository.mu.RUnlock()

	gameNames := make([]string, 0, len(gameRepository.games))

	for _, game := range gameRepository.games {
		gameNames = append(gameNames, game.Name)
	}
	return gameNames, nil
}
