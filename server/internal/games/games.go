package games

import (
	"errors"

	"github.com/kaviraj-j/duoplay/internal/games/tictactoe"
	"github.com/kaviraj-j/duoplay/internal/model"
)

func CreateGameFromName(gameName string) (model.Game, error) {
	switch gameName {
	case "tictactoe":
		return tictactoe.NewTicTacToe(), nil
	}
	return nil, errors.New("game not found for the given name")
}
