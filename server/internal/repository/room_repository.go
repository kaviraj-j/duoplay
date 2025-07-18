package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/kaviraj-j/duoplay/internal/model"
)

type RoomRepository interface {
	CreateRoom(ctx context.Context, room model.Room) error
	GetRoomByID(ctx context.Context, id string) (*model.Room, error)
	AddPlayerToRoom(ctx context.Context, roomID string, player model.Player) error
	DeleteRoom(ctx context.Context, roomID string) error
	UpdateRoom(ctx context.Context, room model.Room) error

	SetGame(ctx context.Context, roomID string, game model.Game) error
	GetGame(ctx context.Context, roomID string) (*model.Game, error)
}

type inMemoryRoomRepository struct {
	rooms map[string]*model.Room
	mu    sync.RWMutex
}

var (
	ErrRoomNotFound error = fmt.Errorf("room not found")
)

func NewRoomRepository() RoomRepository {
	return &inMemoryRoomRepository{
		rooms: make(map[string]*model.Room),
	}
}

func (roomRepository *inMemoryRoomRepository) CreateRoom(ctx context.Context, room model.Room) error {
	roomRepository.mu.Lock()
	defer roomRepository.mu.Unlock()
	roomRepository.rooms[room.ID] = &room
	return nil
}

func (roomRepository *inMemoryRoomRepository) GetRoomByID(ctx context.Context, id string) (*model.Room, error) {
	roomRepository.mu.RLock()
	defer roomRepository.mu.RUnlock()
	room, ok := roomRepository.rooms[id]
	if !ok {
		return nil, ErrRoomNotFound
	}
	return room, nil
}
func (roomRepository *inMemoryRoomRepository) AddPlayerToRoom(ctx context.Context, roomID string, player model.Player) error {
	roomRepository.mu.Lock()
	defer roomRepository.mu.Unlock()
	room, ok := roomRepository.rooms[roomID]
	if !ok {
		return ErrRoomNotFound
	}
	// add user to room
	if len(room.Players) >= 2 {
		return errors.New("player limit exceeded to join room")
	}
	if _, exist := room.Players[player.User.ID]; exist {
		return errors.New("player has already joined room")
	}
	room.Players[player.User.ID] = player
	return nil
}
func (roomRepository *inMemoryRoomRepository) DeleteRoom(ctx context.Context, roomID string) error {
	roomRepository.mu.Lock()
	defer roomRepository.mu.Unlock()
	_, ok := roomRepository.rooms[roomID]
	if !ok {
		return ErrRoomNotFound
	}
	delete(roomRepository.rooms, roomID)
	return nil
}
func (roomRepository *inMemoryRoomRepository) UpdateRoom(ctx context.Context, room model.Room) error {
	roomRepository.mu.Lock()
	defer roomRepository.mu.Unlock()
	existingRoom, ok := roomRepository.rooms[room.ID]
	if !ok {
		return ErrRoomNotFound
	}
	*existingRoom = room
	return nil
}
func (roomRepository *inMemoryRoomRepository) SetGame(ctx context.Context, roomID string, game model.Game) error {
	roomRepository.mu.Lock()
	defer roomRepository.mu.Unlock()
	room, ok := roomRepository.rooms[roomID]
	if !ok {
		return ErrRoomNotFound
	}
	room.Game = game
	return nil
}
func (roomRepository *inMemoryRoomRepository) GetGame(ctx context.Context, roomID string) (*model.Game, error) {
	roomRepository.mu.Lock()
	defer roomRepository.mu.Unlock()
	room, ok := roomRepository.rooms[roomID]
	if !ok {
		return nil, ErrRoomNotFound
	}
	return &room.Game, nil
}
