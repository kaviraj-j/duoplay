package service

import (
	"context"
	"fmt"

	"github.com/kaviraj-j/duoplay/internal/model"
	"github.com/kaviraj-j/duoplay/internal/repository"
)

type RoomService struct {
	roomRepo repository.RoomRepository
}

func NewRoomService(roomRepo repository.RoomRepository) *RoomService {
	return &RoomService{roomRepo: roomRepo}
}

func (s *RoomService) CreateRoom(ctx context.Context) (*model.Room, error) {
	room := model.NewRoom()
	err := s.roomRepo.CreateRoom(ctx, room)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *RoomService) AddPlayer(ctx context.Context, roomID string, player model.Player) error {
	return s.roomRepo.AddPlayerToRoom(ctx, roomID, player)
}

func (s *RoomService) GetRoom(ctx context.Context, roomID string) (*model.Room, error) {
	return s.roomRepo.GetRoomByID(ctx, roomID)
}

func (s *RoomService) StartGame(ctx context.Context, roomID string) error {
	room, err := s.roomRepo.GetRoomByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("error getting room")
	}
	err = room.Game.Start()
	if err != nil {
		return fmt.Errorf("error starting game")
	}
	return nil
}

func (s *RoomService) GetGame(ctx context.Context, roomID string) (*model.Game, error) {
	return s.roomRepo.GetGame(ctx, roomID)
}

// TODO: implement JoinQueue
func (s *RoomService) JoinQueue(ctx context.Context, userID string) error {
	return nil
}
