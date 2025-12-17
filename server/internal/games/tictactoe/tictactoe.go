package tictactoe

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
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
	Board         [3][3]string `json:"Board"`
	CurrentPlayer string       `json:"CurrentPlayer"`
}

type TicTacToe struct {
	state     *model.GameState
	gameState TicTacToeState
}

type Move struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

func NewTicTacToe() model.Game {
	game := &TicTacToe{
		state: &model.GameState{
			ID:        fmt.Sprintf("tictactoe-%d-%v", time.Now().Unix(), getRandomId(100000)),
			Type:      model.TicTacToeGame,
			Name:      "Tic Tac Toe",
			Players:   make(map[string]model.Player),
			Status:    model.GameStatusNotStarted,
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
	t.state.Status = model.GameStatusInProgress
	// Set first player
	for playerID := range t.state.Players {
		t.gameState.CurrentPlayer = playerID
		break
	}
	return nil
}

func (t *TicTacToe) MakeMove(playerID string, move any) error {
	if t.state.Status != model.GameStatusInProgress {
		return fmt.Errorf("game not started")
	}
	if playerID != t.gameState.CurrentPlayer {
		return fmt.Errorf("not your turn")
	}
	var moveData Move
	moveBytes, err := json.Marshal(move)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(moveBytes, &moveData); err != nil {
		return err
	}

	if moveData.Row < 0 || moveData.Row > 2 || moveData.Col < 0 || moveData.Col > 2 {
		return fmt.Errorf("invalid move")
	}

	if t.gameState.Board[moveData.Row][moveData.Col] != "" {
		return fmt.Errorf("position already taken")
	}

	t.gameState.Board[moveData.Row][moveData.Col] = playerID

	// Switch player
	for pid := range t.state.Players {
		if pid != playerID {
			t.gameState.CurrentPlayer = pid
			break
		}
	}

	if t.IsGameOver() {
		t.state.Status = model.GameStatusOver
	}

	return nil
}

func (t *TicTacToe) IsGameOver() bool {
	if t.state.Status == model.GameStatusOver {
		return true
	}

	// Check rows
	for i := 0; i < 3; i++ {
		if t.gameState.Board[i][0] != "" &&
			t.gameState.Board[i][0] == t.gameState.Board[i][1] &&
			t.gameState.Board[i][1] == t.gameState.Board[i][2] {
			// Found a winner
			winnerID := t.gameState.Board[i][0]
			if winner, exists := t.state.Players[winnerID]; exists {
				t.state.Winner = &winner
			}
			t.state.Status = model.GameStatusOver
			return true
		}
	}

	// Check columns
	for i := 0; i < 3; i++ {
		if t.gameState.Board[0][i] != "" &&
			t.gameState.Board[0][i] == t.gameState.Board[1][i] &&
			t.gameState.Board[1][i] == t.gameState.Board[2][i] {
			// Found a winner
			winnerID := t.gameState.Board[0][i]
			if winner, exists := t.state.Players[winnerID]; exists {
				t.state.Winner = &winner
			}
			t.state.Status = model.GameStatusOver
			return true
		}
	}

	// Check diagonals
	if t.gameState.Board[0][0] != "" &&
		t.gameState.Board[0][0] == t.gameState.Board[1][1] &&
		t.gameState.Board[1][1] == t.gameState.Board[2][2] {
		// Found a winner
		winnerID := t.gameState.Board[0][0]
		if winner, exists := t.state.Players[winnerID]; exists {
			t.state.Winner = &winner
		}
		t.state.Status = model.GameStatusOver
		return true
	}
	if t.gameState.Board[0][2] != "" &&
		t.gameState.Board[0][2] == t.gameState.Board[1][1] &&
		t.gameState.Board[1][1] == t.gameState.Board[2][0] {
		// Found a winner
		winnerID := t.gameState.Board[0][2]
		if winner, exists := t.state.Players[winnerID]; exists {
			t.state.Winner = &winner
		}
		t.state.Status = model.GameStatusOver
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
	// It's a draw - no winner
	t.state.Winner = nil
	t.state.Status = model.GameStatusOver
	return true
}

func (t *TicTacToe) ResetState() error {
	t.gameState = TicTacToeState{
		Board: [3][3]string{
			{"", "", ""},
			{"", "", ""},
			{"", "", ""},
		},
	}
	return nil
}

func (t *TicTacToe) GetWinner() *model.Player {
	return t.state.Winner
}

func (t *TicTacToe) GetStatus() model.GameStatus {
	return t.state.Status
}

// SetPlayers sets the players for the game
func (t *TicTacToe) SetPlayers(players map[string]model.Player) {
	t.state.Players = players
}
