package model

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// MessageType is the type for WebSocket message types
// Use these constants for message type enums in both server and client
// You can export these to the client via codegen or documentation

type MessageType string

const (
	MessageTypeJoinRoom        MessageType = "join_room"
	MessageTypeOpponentJoined  MessageType = "opponent_joined"
	MessageTypeChooseGame      MessageType = "choose_game"
	MessageTypeGameMove        MessageType = "game_move"
	MessageTypeGameChosen      MessageType = "game_chosen"
	MessageTypeGameAccepted    MessageType = "game_accepted"
	MessageTypeMoveMade        MessageType = "move_made"
	MessageTypeRejectGame      MessageType = "reject_game"
	MessageTypeContinueGame    MessageType = "continue_game"
	MessageTypeError           MessageType = "error"
	MessageTypeGameSelection   MessageType = "game_selection"
	MessageTypeBothGamesChosen MessageType = "both_games_chosen"
)

type RoomStatus string

const (
	RoomStatusWaitingForPlayer RoomStatus = "waiting_for_player"
	RoomStatusGameSelection    RoomStatus = "game_selection"
	RoomStatusGameSelected     RoomStatus = "game_selected"
	RoomStatusGameStarted      RoomStatus = "game_started"
	RoomStatusGameOver         RoomStatus = "game_over"
)

type Player struct {
	User User
	Conn *websocket.Conn
}

// GameSelectionState tracks which players have chosen games
type GameSelectionState struct {
	PlayerChoices map[string]GameType `json:"player_choices"` // playerID -> gameType
	BothChosen    bool                `json:"both_chosen"`
}

type Room struct {
	ID            string             `json:"id"`
	Players       map[string]Player  `json:"players"`
	Game          Game               `json:"game"`
	GameSelection GameSelectionState `json:"game_selection"`
	IsGameStarted bool               `json:"is_game_started"`
	Status        RoomStatus         `json:"status"`
}

func NewRoom() Room {
	return Room{
		ID:      uuid.New().String(),
		Players: make(map[string]Player),
		GameSelection: GameSelectionState{
			PlayerChoices: make(map[string]GameType),
			BothChosen:    false,
		},
		IsGameStarted: false,
		Status:        RoomStatusWaitingForPlayer,
	}
}
