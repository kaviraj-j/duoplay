package service

import (
	"context"

	"github.com/kaviraj-j/duoplay/internal/model"
	"github.com/kaviraj-j/duoplay/internal/repository"
)

type GameService struct {
	gameRepo repository.GameRepository
}

func NewGameService(gameRepo repository.GameRepository) *GameService {
	return &GameService{gameRepo: gameRepo}
}

func (s *GameService) GetGamesList(ctx context.Context) ([]model.GameListPayload, error) {
	return s.gameRepo.GetGamesList(ctx)
}
