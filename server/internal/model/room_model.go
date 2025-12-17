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
	MessageTypeGameAccept      MessageType = "game_accept"
	MessageTypeGameAccepted    MessageType = "game_accepted"
	MessageTypeGameReject      MessageType = "game_reject"
	MessageTypeGameRejected    MessageType = "game_rejected"
	MessageTypeStartGame       MessageType = "start_game"
	MessageTypeMoveMade        MessageType = "move_made"
	MessageTypeReplayGame      MessageType = "replay_game"
	MessageTypeReplayAccepted  MessageType = "replay_accepted"
	MessageTypeReplayRejected  MessageType = "replay_rejected"
	MessageTypeError           MessageType = "error"
	MessageTypeGameSelection   MessageType = "game_selection"
	MessageTypeBothGamesChosen MessageType = "both_games_chosen"
	MessageTypeAuth            MessageType = "auth"
	MessageTypeRoomCreated     MessageType = "room_created"
	MessageTypeQueueJoined     MessageType = "queue_joined"
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
}

type RoomEventType string

const (
	RoomEventTypeGameOver RoomEventType = "game_over"
)

type Event struct {
	Type    RoomEventType `json:"type"`
	Payload interface{}   `json:"payload"`
}

type Room struct {
	ID            string             `json:"id"`
	Players       map[string]Player  `json:"players"`
	Game          Game               `json:"game"`
	GameSelection GameSelectionState `json:"game_selection"`
	Status        RoomStatus         `json:"status"`
	EventChannel  chan Event         `json:"event_channel"`
}

type RoomPlayer struct {
	User User `json:"user"`
}

type RoomResponse struct {
	ID            string                `json:"id"`
	Players       map[string]RoomPlayer `json:"players"`
	Status        RoomStatus            `json:"status"`
	GameSelection map[string]GameType   `json:"game_selection"`
	Game          *GameResponse         `json:"game,omitempty"`
}

type GameResponse struct {
	Type   GameType    `json:"type"`
	Status GameStatus  `json:"status"`
	State  interface{} `json:"state"`
}

func NewRoom() Room {
	return Room{
		ID:      uuid.New().String(),
		Players: make(map[string]Player),
		GameSelection: GameSelectionState{
			PlayerChoices: make(map[string]GameType),
		},
		Status:       RoomStatusWaitingForPlayer,
		EventChannel: make(chan Event),
	}
}

func (r Room) GetRoomResponse() RoomResponse {
	players := make(map[string]RoomPlayer)
	for id, p := range r.Players {
		players[id] = RoomPlayer{
			User: p.User,
		}
	}

	response := RoomResponse{
		ID:            r.ID,
		Players:       players,
		Status:        r.Status,
		GameSelection: r.GameSelection.PlayerChoices,
	}

	// Include game information if game exists
	if r.Game != nil {
		response.Game = &GameResponse{
			Type:   r.Game.GetType(),
			Status: r.Game.GetStatus(),
			State:  r.Game.GetState(),
		}
	}

	return response
}
