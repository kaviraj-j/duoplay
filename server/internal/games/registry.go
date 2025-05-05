package games

import (
	"fmt"

	"github.com/kaviraj-j/duoplay/internal/model"
)

type GameFactory interface {
	CreateGame(state *model.GameState) (model.Game, error)
}

type gameRegistry struct {
	factories map[model.GameType]GameFactory
}

func NewGameRegistry() *gameRegistry {
	return &gameRegistry{
		factories: make(map[model.GameType]GameFactory),
	}
}

func (r *gameRegistry) Register(gameType model.GameType, factory GameFactory) {
	r.factories[gameType] = factory
}

func (r *gameRegistry) CreateGame(gameType model.GameType, state *model.GameState) (model.Game, error) {
	factory, exists := r.factories[gameType]
	if !exists {
		return nil, fmt.Errorf("no factory registered for game type: %s", gameType)
	}
	return factory.CreateGame(state)
}
