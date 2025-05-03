package model

import (
	"crypto/rand"
	"math/big"

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

const idLength = 12
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomID() string {
	result := make([]byte, idLength)
	for i := 0; i < idLength; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

func NewRoom() Room {
	return Room{
		ID:      generateRandomID(),
		Players: make(map[string]Player),
	}
}
