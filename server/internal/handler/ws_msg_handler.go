package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kaviraj-j/duoplay/internal/model"
)

// handleWebSocketMessages handles WebSocket communication for a game session
func (h *RoomHandler) handleWebSocketMessages(c *gin.Context, conn *websocket.Conn, roomID string, player model.Player) {
	defer conn.Close()

	for {
		// Read message from WebSocket
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Parse message as JSON
		var msg map[string]interface{}
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			conn.WriteJSON(WSMessage{Type: "error", Message: "Invalid message format", Data: nil})
			continue
		}
		typeStr, ok := msg["type"].(string)
		if !ok {
			conn.WriteJSON(WSMessage{Type: "error", Message: "Missing or invalid message type", Data: nil})
			continue
		}

		typeVal := model.MessageType(typeStr)
		switch typeVal {
		// when a player joins a room
		case model.MessageTypeJoinRoom:
			h.handlePlayerJoinedRoom(c, conn, roomID, player)
		case model.MessageTypeChooseGame:
			h.handleGameChosen(c, conn, roomID, player, msg)
		case model.MessageTypeGameAccept:
			h.handleGameAccepted(c, conn, roomID, player, msg)
		case model.MessageTypeGameReject:
			h.handleGameRejected(c, conn, roomID, player, msg)
		case model.MessageTypeGameMove:
			h.handleGameMove(c, conn, roomID, player, msg)
		case model.MessageTypeReplayGame:
			h.handleReplayGame(c, conn, roomID, player, msg)
		case model.MessageTypeReplayAccepted:
			h.handleReplayAccepted(c, conn, roomID, player, msg)
		case model.MessageTypeReplayRejected:
			h.handleReplayRejected(c, conn, roomID, player, msg)
		default:
			conn.WriteJSON(WSMessage{Type: "error", Message: "Unknown message type", Data: nil})
		}
	}
}

// handlePlayerJoinedRoom handles the player joining a room
func (h *RoomHandler) handlePlayerJoinedRoom(c *gin.Context, conn *websocket.Conn, roomID string, player model.Player) {
	// check if the room is valid
	room, err := h.roomService.GetRoom(c, roomID)
	if err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Room not found", Data: nil})
		return
	}

	// check if the player is in the room
	_, exists := room.Players[player.User.ID]
	if !exists {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Player not found in room", Data: nil})
		return
	}

	// Use the service method to handle player joined room
	if err := h.roomService.HandlePlayerJoinedRoom(c, room, player); err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Failed to notify other players", Data: nil})
		return
	}

	// Send confirmation to the joining player
	conn.WriteJSON(WSMessage{Type: "joined_room", Message: "You joined the room", Data: nil})
}

// handleGameChosen handles a player choosing a game
func (h *RoomHandler) handleGameChosen(c *gin.Context, conn *websocket.Conn, roomID string, player model.Player, msg map[string]interface{}) {
	// check if the room is valid
	room, err := h.roomService.GetRoom(c, roomID)
	if err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Room not found", Data: nil})
		return
	}

	// Get the game type from the parsed message
	gameTypeStr, ok := msg["game_type"].(string)
	if !ok {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Game type is required", Data: nil})
		return
	}

	gameType := model.GameType(gameTypeStr)

	// Handle the game choice through the service
	if err := h.roomService.HandleGameChosen(c, room, player, gameType); err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Failed to handle game choice", Data: nil})
		return
	}

	// Send confirmation to the player who chose the game
	conn.WriteJSON(WSMessage{
		Type:    "game_chosen_confirmation",
		Message: "Your game choice has been recorded",
		Data: map[string]interface{}{
			"game_type": gameType,
		},
	})
}

// handleGameAccepted handles a player accepting a game
func (h *RoomHandler) handleGameAccepted(c *gin.Context, conn *websocket.Conn, roomID string, player model.Player, msg map[string]interface{}) {
	// check if the room is valid
	room, err := h.roomService.GetRoom(c, roomID)
	if err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Room not found", Data: nil})
		return
	}

	// Get the game type from the parsed message
	gameTypeStr, ok := msg["game_type"].(string)
	if !ok {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Game type is required", Data: nil})
		return
	}

	gameType := model.GameType(gameTypeStr)

	// Handle the game choice through the service
	if err := h.roomService.HandleGameAccepted(c, room, player, gameType); err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Failed to handle game accepted", Data: nil})
		return
	}

}

// handleGameRejected handles a player rejecting a game
func (h *RoomHandler) handleGameRejected(c *gin.Context, conn *websocket.Conn, roomID string, player model.Player, msg map[string]interface{}) {
	// check if the room is valid
	room, err := h.roomService.GetRoom(c, roomID)
	if err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Room not found", Data: nil})
		return
	}

	// Get the game type from the parsed message
	gameTypeStr, ok := msg["game_type"].(string)
	if !ok {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Game type is required", Data: nil})
		return
	}

	gameType := model.GameType(gameTypeStr)

	// Handle the game choice through the service
	if err := h.roomService.HandleGameRejected(c, room, player, gameType); err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Failed to handle game rejected", Data: nil})
		return
	}
}

// handleGameMove handles a player making a move in the game
func (h *RoomHandler) handleGameMove(c *gin.Context, conn *websocket.Conn, roomID string, player model.Player, msg map[string]interface{}) {
	room, err := h.roomService.GetRoom(c, roomID)
	if err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Room not found", Data: nil})
		return
	}

	_, exists := room.Players[player.User.ID]
	if !exists {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Player not found in room", Data: nil})
		return
	}

	if room.Game == nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Game not started", Data: nil})
		return
	}

	move, ok := msg["move"]
	if !ok {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Move is required", Data: nil})
		return
	}

	// Make the move
	if err := room.Game.MakeMove(player.User.ID, move); err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: err.Error(), Data: nil})
		return
	}

	// Update room status if game is over
	if room.Game.IsGameOver() {
		room.Status = model.RoomStatusGameOver
		roomResponse := room.GetRoomResponse()
		room.EventChannel <- model.Event{
			Type:    model.RoomEventTypeGameOver,
			Payload: roomResponse,
		}
	}

	// Update room in repository
	if err := h.roomService.UpdateRoom(c, *room); err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Failed to update room", Data: nil})
		return
	}

	// Get updated room response with game state
	roomResponse := room.GetRoomResponse()

	// Broadcast the move to both players
	for _, p := range room.Players {
		if p.Conn != nil {
			p.Conn.WriteJSON(WSMessage{
				Type:    string(model.MessageTypeMoveMade),
				Message: "Move made",
				Data:    roomResponse,
			})
		}
	}
}

func (h *RoomHandler) handleReplayGame(c *gin.Context, conn *websocket.Conn, roomID string, player model.Player, msg map[string]interface{}) {
	room, err := h.roomService.GetRoom(c, roomID)
	if err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Room not found", Data: nil})
		return
	}

	_, exists := room.Players[player.User.ID]
	if !exists {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Player not found in room", Data: nil})
		return
	}

	if err := h.roomService.HandleReplayGame(c, room, player, msg); err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Failed to handle replay game", Data: nil})
		return
	}

	conn.WriteJSON(WSMessage{
		Type:    "replay_game_received",
		Message: "Replay game received, processing...",
		Data:    nil,
	})
}

func (h *RoomHandler) handleReplayAccepted(c *gin.Context, conn *websocket.Conn, roomID string, player model.Player, msg map[string]interface{}) {
	room, err := h.roomService.GetRoom(c, roomID)
	if err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Room not found", Data: nil})
		return
	}

	if err := h.roomService.HandleReplayAccepted(c, room, player, msg); err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Failed to handle replay accepted", Data: nil})
		return
	}
}

func (h *RoomHandler) handleReplayRejected(c *gin.Context, conn *websocket.Conn, roomID string, player model.Player, msg map[string]interface{}) {
	room, err := h.roomService.GetRoom(c, roomID)
	if err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Room not found", Data: nil})
		return
	}

	if err := h.roomService.HandleReplayRejected(c, room, player, msg); err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Failed to handle replay rejected", Data: nil})
		return
	}
}
