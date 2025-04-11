package model

import (
	"golang.org/x/net/websocket"
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
