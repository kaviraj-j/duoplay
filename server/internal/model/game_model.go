package model

import "time"

type GameType string

const (
	TicTacToeGame GameType = "tictactoe"
	// TODO: add new games to the list
)

// Game interface defines core game behavior
type Game interface {
	GetType() GameType
	GetState() any
	Start() error
	MakeMove(playerID string, move any) error
	IsGameOver() bool
}

// GameState holds common game state
type GameState struct {
	ID        string
	Type      GameType
	Name      string
	Players   map[string]Player
	Winner    *Player
	IsStarted bool
	CreatedAt time.Time
}

type GameListPayload struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}
