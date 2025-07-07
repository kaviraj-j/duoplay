package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kaviraj-j/duoplay/internal/games/tictactoe"
	"github.com/kaviraj-j/duoplay/internal/model"
	"github.com/kaviraj-j/duoplay/internal/repository"
)

var ErrorUserAlreadyInQueue = errors.New("user is already in queue")

type RoomService struct {
	roomRepo  repository.RoomRepository
	userRepo  repository.UserRepository
	queueRepo repository.QueueRepository
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewRoomService(roomRepo repository.RoomRepository, queueRepo repository.QueueRepository, userRepo repository.UserRepository) *RoomService {
	ctx, cancel := context.WithCancel(context.Background())
	service := &RoomService{
		roomRepo:  roomRepo,
		userRepo:  userRepo,
		queueRepo: queueRepo,
		ctx:       ctx,
		cancel:    cancel,
	}

	// Set up the match callback
	queueRepo.SetMatchCallback(service.CreateMatch)

	// Start the centralized queue monitor
	go service.monitorQueue(ctx)

	return service
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

func (s *RoomService) JoinQueue(ctx context.Context, userID string, conn *websocket.Conn) error {
	// check if user is already in queue
	isInQueue, err := s.queueRepo.PlayerExistsInQueue(ctx, userID)
	if err != nil {
		return fmt.Errorf("error checking if user is in queue")
	}

	if isInQueue {
		return ErrorUserAlreadyInQueue
	}

	// Add player to queue
	err = s.queueRepo.AddToQueue(ctx, userID, conn)
	if err != nil {
		return err
	}
	return nil
}

func (s *RoomService) RemoveFromQueue(ctx context.Context, userID string) error {
	return s.queueRepo.RemoveFromQueue(ctx, userID)
}

func (s *RoomService) monitorQueue(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second) // Check every 2 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Get all waiting players
			waitingPlayers, err := s.queueRepo.GetWaitingPlayers(ctx)
			if err != nil {
				continue
			}

			// If we have 2 or more players, create matches
			for len(waitingPlayers) >= 2 {
				player1ID := waitingPlayers[0]
				player2ID := waitingPlayers[1]

				// Try to create a match
				room, err := s.CreateMatch(ctx, player1ID, player2ID)
				if err != nil {
					// If match creation fails, remove problematic players from queue
					s.queueRepo.RemoveFromQueue(ctx, player1ID)
					s.queueRepo.RemoveFromQueue(ctx, player2ID)
					continue
				}

				// Remove matched players from waiting list
				waitingPlayers = waitingPlayers[2:]

				// Log successful match
				fmt.Printf("Match created: Room %s with players %s and %s\n", room.ID, player1ID, player2ID)
			}
		}
	}
}

func (s *RoomService) CreateMatch(ctx context.Context, player1ID, player2ID string) (*model.Room, error) {
	// Create a new room
	room := model.NewRoom()

	// Get player connections from queue
	player1Conn, err := s.queueRepo.GetPlayerConnection(ctx, player1ID)
	if err != nil {
		return nil, fmt.Errorf("player 1 connection not found: %v", err)
	}

	player2Conn, err := s.queueRepo.GetPlayerConnection(ctx, player2ID)
	if err != nil {
		return nil, fmt.Errorf("player 2 connection not found: %v", err)
	}

	// Create player objects
	user1, err := s.userRepo.FindByID(ctx, player1ID)
	if err != nil {
		return nil, fmt.Errorf("user 1 not found: %v", err)
	}
	player1 := model.Player{
		User: *user1,
		Conn: player1Conn,
	}

	user2, err := s.userRepo.FindByID(ctx, player2ID)
	if err != nil {
		return nil, fmt.Errorf("user 2 not found: %v", err)
	}
	player2 := model.Player{
		User: *user2,
		Conn: player2Conn,
	}

	// Add players to room
	room.Players[player1ID] = player1
	room.Players[player2ID] = player2

	// Create a game instance
	game := tictactoe.NewTicTacToe()
	room.Game = game

	// Save room to repository
	err = s.roomRepo.CreateRoom(ctx, room)
	if err != nil {
		return nil, err
	}

	// Notify both players about the match
	player1Conn.WriteJSON(gin.H{
		"type":    "match_found",
		"room_id": room.ID,
		"message": "Match found! Game starting...",
	})

	player2Conn.WriteJSON(gin.H{
		"type":    "match_found",
		"room_id": room.ID,
		"message": "Match found! Game starting...",
	})

	return &room, nil
}

func (s *RoomService) LeaveRoom(ctx context.Context, roomID string) error {
	_, err := s.roomRepo.GetRoomByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("error getting room")
	}

	return s.roomRepo.DeleteRoom(ctx, roomID)
}
