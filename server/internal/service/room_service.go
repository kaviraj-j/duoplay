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

	// Start the centralized queue monitor
	go service.monitorQueue(ctx)

	return service
}

func (s *RoomService) CreateRoom(ctx context.Context) (model.Room, error) {
	room := model.NewRoom()
	err := s.roomRepo.CreateRoom(ctx, room)
	if err != nil {
		return model.Room{}, err
	}
	return room, nil
}

func (s *RoomService) AddPlayer(ctx context.Context, roomID string, player model.Player) error {
	return s.roomRepo.AddPlayerToRoom(ctx, roomID, player)
}

func (s *RoomService) GetRoom(ctx context.Context, roomID string) (*model.Room, error) {
	return s.roomRepo.GetRoomByID(ctx, roomID)
}

func (s *RoomService) UpdateRoom(ctx context.Context, room model.Room) error {
	return s.roomRepo.UpdateRoom(ctx, room)
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

	// update the room status to game started
	room.Status = model.RoomStatusGameStarted
	err = s.roomRepo.UpdateRoom(ctx, *room)
	if err != nil {
		return err
	}
	return nil
}

func (s *RoomService) GetGame(ctx context.Context, roomID string) (*model.Game, error) {
	return s.roomRepo.GetGame(ctx, roomID)
}

func (s *RoomService) JoinQueue(ctx context.Context, userID string, conn *websocket.Conn) error {
	// check if user is already in queue
	isInQueue := s.queueRepo.PlayerExistsInQueue(ctx, userID)

	if isInQueue {
		return ErrorUserAlreadyInQueue
	}

	// Add player to queue
	err := s.queueRepo.AddToQueue(ctx, userID, conn)
	if err != nil {
		return err
	}
	return nil
}

func (s *RoomService) RemoveFromQueue(ctx context.Context, userID string) error {
	return s.queueRepo.RemoveFromQueue(ctx, userID)
}

func (s *RoomService) monitorQueue(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond)
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
					continue
				}

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

	// Update room status to game selection
	room.Status = model.RoomStatusGameSelection
	err = s.roomRepo.UpdateRoom(ctx, room)
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

func (s *RoomService) GetOppositePlayer(ctx context.Context, roomPlayers map[string]model.Player, playerId string) (model.Player, error) {
	for id, p := range roomPlayers {
		if id != playerId {
			return p, nil
		}
	}
	return model.Player{}, errors.New("opposite player not found")
}

// handlePlayerJoinedRoom notifies the other player in the room that a player has joined
func (s *RoomService) HandlePlayerJoinedRoom(ctx context.Context, room *model.Room, joinedPlayer model.Player) error {
	// Find the opposite player
	oppositePlayer, err := s.GetOppositePlayer(ctx, room.Players, joinedPlayer.User.ID)
	if err != nil {
		return err
	}
	if oppositePlayer.Conn != nil {
		err := oppositePlayer.Conn.WriteJSON(map[string]interface{}{
			"type":    "joined_room",
			"message": "A player has joined your room.",
			"data": map[string]interface{}{
				"user": joinedPlayer.User,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// HandleGameChosen handles when a player chooses a game
func (s *RoomService) HandleGameChosen(ctx context.Context, room *model.Room, player model.Player, gameType model.GameType) error {
	// Record the player's game choice
	room.GameSelection.PlayerChoices[player.User.ID] = gameType

	// Notify the opposite player about the game choice
	oppositePlayer, err := s.GetOppositePlayer(ctx, room.Players, player.User.ID)
	if err != nil {
		return err
	}

	if oppositePlayer.Conn != nil {
		messageData := map[string]interface{}{
			"player_id":   player.User.ID,
			"player_name": player.User.Name,
			"game_type":   gameType,
		}

		err := oppositePlayer.Conn.WriteJSON(map[string]interface{}{
			"type":    model.MessageTypeGameChosen,
			"message": "Your opponent has chosen a game.",
			"data":    messageData,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// GetGameSelectionState returns the current game selection state for a room
func (s *RoomService) GetGameSelectionState(ctx context.Context, roomID string) (*model.GameSelectionState, error) {
	room, err := s.roomRepo.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	return &room.GameSelection, nil
}

// HandleGameAccepted handles when a player accepts a game
func (s *RoomService) HandleGameAccepted(ctx context.Context, room *model.Room, player model.Player, gameType model.GameType) error {

	// check if the opposite player has choose this game
	oppositePlayer, err := s.GetOppositePlayer(ctx, room.Players, player.User.ID)
	if err != nil {
		return err
	}

	if room.GameSelection.PlayerChoices[oppositePlayer.User.ID] != gameType {
		return errors.New("opposite player has not chosen this game")
	}

	// Record the player's game acceptance
	room.GameSelection.PlayerChoices[player.User.ID] = gameType

	// notify the opposite player about the game acceptance
	if oppositePlayer.Conn != nil {
		err = oppositePlayer.Conn.WriteJSON(map[string]interface{}{
			"type":    model.MessageTypeGameAccepted,
			"message": "Your opponent has accepted the game.",
		})
		if err != nil {
			return err
		}
	}

	// Send start_game message to both players when game is accepted
	if player.Conn != nil {
		err = player.Conn.WriteJSON(map[string]interface{}{
			"type":      model.MessageTypeStartGame,
			"game_type": gameType,
		})
		if err != nil {
			return err
		}
	}

	if oppositePlayer.Conn != nil {
		err = oppositePlayer.Conn.WriteJSON(map[string]interface{}{
			"type":      model.MessageTypeStartGame,
			"game_type": gameType,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// HandleGameRejected handles when a player rejects a game
func (s *RoomService) HandleGameRejected(ctx context.Context, room *model.Room, player model.Player, gameType model.GameType) error {
	// check if the opposite player has choose this game
	oppositePlayer, err := s.GetOppositePlayer(ctx, room.Players, player.User.ID)
	if err != nil {
		return err
	}

	// notify the opposite player about the game rejection
	if oppositePlayer.Conn != nil {
		err = oppositePlayer.Conn.WriteJSON(map[string]interface{}{
			"type":    model.MessageTypeGameRejected,
			"message": "Your opponent has rejected the game.",
		})
	}

	return nil
}

func (s *RoomService) HandleReplayGame(ctx context.Context, room *model.Room, player model.Player, msg map[string]interface{}) error {
	// check if the opposite player has choose this game
	oppositePlayer, err := s.GetOppositePlayer(ctx, room.Players, player.User.ID)
	if err != nil {
		return err
	}

	// notify the opposite player about the replay game
	if oppositePlayer.Conn != nil {
		err = oppositePlayer.Conn.WriteJSON(map[string]interface{}{
			"type":    model.MessageTypeReplayGame,
			"message": "Your opponent has requested a replay.",
		})
	}

	return nil
}

func (s *RoomService) HandleReplayAccepted(ctx context.Context, room *model.Room, player model.Player, msg map[string]interface{}) error {
	// check if the opposite player has choose this game
	oppositePlayer, err := s.GetOppositePlayer(ctx, room.Players, player.User.ID)
	if err != nil {
		return err
	}

	// notify the opposite player about the replay accepted
	if oppositePlayer.Conn != nil {
		err = oppositePlayer.Conn.WriteJSON(map[string]interface{}{
			"type":    model.MessageTypeReplayAccepted,
			"message": "Your opponent has accepted the replay.",
		})
	}

	// reset the game state
	room.Game.ResetState()

	s.StartGame(ctx, room.ID)

	return nil
}

func (s *RoomService) HandleReplayRejected(ctx context.Context, room *model.Room, player model.Player, msg map[string]interface{}) error {
	// check if the opposite player has choose this game
	oppositePlayer, err := s.GetOppositePlayer(ctx, room.Players, player.User.ID)
	if err != nil {
		return err
	}

	// notify the opposite player about the replay rejected
	if oppositePlayer.Conn != nil {
		err = oppositePlayer.Conn.WriteJSON(map[string]interface{}{
			"type":    model.MessageTypeReplayRejected,
			"message": "Your opponent has rejected the replay.",
		})
	}

	return nil
}
