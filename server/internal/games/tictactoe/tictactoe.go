package tictactoe

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/kaviraj-j/duoplay/internal/model"
)

func getRandomId(max int) string {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d", int(binary.BigEndian.Uint64(b[:]))%max)
}

type TicTacToeState struct {
	Board         [3][3]string
	CurrentPlayer string
}

type TicTacToe struct {
	state     *model.GameState
	gameState TicTacToeState
}

func NewTicTacToe() model.Game {
	game := &TicTacToe{
		state: &model.GameState{
			ID:        fmt.Sprintf("tictactoe-%d-%v", time.Now().Unix(), getRandomId(100000)),
			Type:      model.TicTacToeGame,
			Name:      "Tic Tac Toe",
			Players:   make(map[string]model.Player),
			IsStarted: false,
			CreatedAt: time.Now(),
		},
		gameState: TicTacToeState{
			Board: [3][3]string{
				{"", "", ""},
				{"", "", ""},
				{"", "", ""},
			},
		},
	}
	return game
}

func (t *TicTacToe) GetType() model.GameType {
	return model.TicTacToeGame
}

func (t *TicTacToe) GetState() any {
	return t.gameState
}

func (t *TicTacToe) Start() error {
	if len(t.state.Players) != 2 {
		return fmt.Errorf("need exactly 2 players to start")
	}
	t.state.IsStarted = true
	// Set first player
	for playerID := range t.state.Players {
		t.gameState.CurrentPlayer = playerID
		break
	}
	return nil
}

func (t *TicTacToe) MakeMove(playerID string, move any) error {
	if !t.state.IsStarted {
		return fmt.Errorf("game not started")
	}
	if playerID != t.gameState.CurrentPlayer {
		return fmt.Errorf("not your turn")
	}

	moveData, ok := move.(map[string]int)
	if !ok {
		return fmt.Errorf("invalid move format")
	}

	row, col := moveData["row"], moveData["col"]
	if row < 0 || row > 2 || col < 0 || col > 2 {
		return fmt.Errorf("invalid position")
	}
	if t.gameState.Board[row][col] != "" {
		return fmt.Errorf("position already taken")
	}

	// Make the move
	t.gameState.Board[row][col] = playerID

	// Switch player
	for pid := range t.state.Players {
		if pid != playerID {
			t.gameState.CurrentPlayer = pid
			break
		}
	}

	return nil
}

func (t *TicTacToe) IsGameOver() bool {
	// Check rows
	for i := 0; i < 3; i++ {
		if t.gameState.Board[i][0] != "" &&
			t.gameState.Board[i][0] == t.gameState.Board[i][1] &&
			t.gameState.Board[i][1] == t.gameState.Board[i][2] {
			return true
		}
	}

	// Check columns
	for i := 0; i < 3; i++ {
		if t.gameState.Board[0][i] != "" &&
			t.gameState.Board[0][i] == t.gameState.Board[1][i] &&
			t.gameState.Board[1][i] == t.gameState.Board[2][i] {
			return true
		}
	}

	// Check diagonals
	if t.gameState.Board[0][0] != "" &&
		t.gameState.Board[0][0] == t.gameState.Board[1][1] &&
		t.gameState.Board[1][1] == t.gameState.Board[2][2] {
		return true
	}
	if t.gameState.Board[0][2] != "" &&
		t.gameState.Board[0][2] == t.gameState.Board[1][1] &&
		t.gameState.Board[1][1] == t.gameState.Board[2][0] {
		return true
	}

	// Check for draw
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if t.gameState.Board[i][j] == "" {
				return false
			}
		}
	}
	return true
}
