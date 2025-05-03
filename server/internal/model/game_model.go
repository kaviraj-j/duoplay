package model

type Game interface {
	Start() error
	IsGameOver() bool
	UpdateState(msg any) error
}
