package service

import (
	"context"

	"github.com/kaviraj-j/duoplay/internal/repository"
)

type GameService struct {
	gameRepo repository.GameRepository
}

func NewGameService(gameRepo repository.GameRepository) *GameService {
	return &GameService{gameRepo: gameRepo}
}

func (s *GameService) GetGamesList(ctx context.Context) ([]string, error) {
	return s.gameRepo.GetGamesList(ctx)
}
