package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/kaviraj-j/duoplay/internal/games/tictactoe"
	"github.com/kaviraj-j/duoplay/internal/model"
)

type GameRepository interface {
	GetGamesList(ctx context.Context) ([]string, error)
	GetGame(ctx context.Context, gameID string) (model.Game, error)
	CreateGame(ctx context.Context, gameType model.GameType) (model.Game, error)
	UpdateGame(ctx context.Context, gameID string, game model.Game) error
}

type inMemoryGameRepository struct {
	availableGames map[model.GameType]struct{}
	activeGames    map[string]model.Game
	mu             sync.RWMutex
}

var (
	ErrGameNotFound error = fmt.Errorf("game not found")
)

func NewGameRepository() GameRepository {
	gameRepo := &inMemoryGameRepository{
		availableGames: map[model.GameType]struct{}{
			model.TicTacToeGame: {}, // Register tic-tac-toe as an available game type
		},
		activeGames: make(map[string]model.Game),
	}
	return gameRepo
}

func (gameRepository *inMemoryGameRepository) GetGamesList(ctx context.Context) ([]string, error) {
	gameRepository.mu.RLock()
	defer gameRepository.mu.RUnlock()

	gameTypes := make([]string, 0, len(gameRepository.availableGames))
	for gameType := range gameRepository.availableGames {
		gameTypes = append(gameTypes, string(gameType))
	}
	return gameTypes, nil
}

func (gameRepository *inMemoryGameRepository) GetGame(ctx context.Context, gameID string) (model.Game, error) {
	gameRepository.mu.RLock()
	defer gameRepository.mu.RUnlock()

	game, exists := gameRepository.activeGames[gameID]
	if !exists {
		return nil, ErrGameNotFound
	}
	return game, nil
}

func (gameRepository *inMemoryGameRepository) CreateGame(ctx context.Context, gameType model.GameType) (model.Game, error) {
	gameRepository.mu.Lock()
	defer gameRepository.mu.Unlock()

	// Check if the game type is available
	if _, exists := gameRepository.availableGames[gameType]; !exists {
		return nil, fmt.Errorf("unsupported game type: %s", gameType)
	}

	var game model.Game
	switch gameType {
	case model.TicTacToeGame:
		game = tictactoe.NewTicTacToe()
	default:
		return nil, fmt.Errorf("unsupported game type: %s", gameType)
	}

	gameID := fmt.Sprintf("%s-%d", gameType, len(gameRepository.activeGames))
	gameRepository.activeGames[gameID] = game
	return game, nil
}

func (gameRepository *inMemoryGameRepository) UpdateGame(ctx context.Context, gameID string, game model.Game) error {
	gameRepository.mu.Lock()
	defer gameRepository.mu.Unlock()

	if _, exists := gameRepository.activeGames[gameID]; !exists {
		return ErrGameNotFound
	}

	gameRepository.activeGames[gameID] = game
	return nil
}
