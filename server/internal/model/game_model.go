package model

import "time"

type GameStatus string

const (
	GameStatusNotStarted GameStatus = "not_started"
	GameStatusInProgress GameStatus = "in_progress"
	GameStatusOver       GameStatus = "over"
)

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
	GetWinner() *Player
	GetStatus() GameStatus
}

// GameState holds common game state
type GameState struct {
	ID        string
	Type      GameType
	Name      string
	Players   map[string]Player
	Winner    *Player
	Status    GameStatus
	CreatedAt time.Time
}

type GameListPayload struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}
