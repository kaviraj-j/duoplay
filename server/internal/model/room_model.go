package model

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	User User
	Conn *websocket.Conn
}

type Room struct {
	ID      string `json:"id"`
	Players map[string]Player
	Game    Game
}

func NewRoom() Room {
	return Room{
		ID:      uuid.New().String(),
		Players: make(map[string]Player),
	}
}
