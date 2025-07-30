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
			conn.WriteJSON(gin.H{"type": "error", "message": "Invalid message format", "data": nil})
			continue
		}

		typeVal, ok := msg["type"].(model.MessageType)
		if !ok {
			conn.WriteJSON(gin.H{"type": "error", "message": "Missing message type", "data": nil})
			continue
		}

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
		default:
			conn.WriteJSON(gin.H{"type": "error", "message": "Unknown message type", "data": nil})
		}
	}
}

// handlePlayerJoinedRoom handles the player joining a room
func (h *RoomHandler) handlePlayerJoinedRoom(c *gin.Context, conn *websocket.Conn, roomID string, player model.Player) {
	// check if the room is valid
	room, err := h.roomService.GetRoom(c, roomID)
	if err != nil {
		conn.WriteJSON(gin.H{"type": "error", "message": "Room not found", "data": nil})
		return
	}

	// check if the player is in the room
	_, exists := room.Players[player.User.ID]
	if !exists {
		conn.WriteJSON(gin.H{"type": "error", "message": "Player not found in room", "data": nil})
		return
	}

	// Use the service method to handle player joined room
	if err := h.roomService.HandlePlayerJoinedRoom(c, room, player); err != nil {
		conn.WriteJSON(gin.H{"type": "error", "message": "Failed to notify other players", "data": nil})
		return
	}

	// Send confirmation to the joining player
	conn.WriteJSON(gin.H{"type": "joined_room", "message": "You joined the room", "data": nil})
}

// handleGameChosen handles a player choosing a game
func (h *RoomHandler) handleGameChosen(c *gin.Context, conn *websocket.Conn, roomID string, player model.Player, msg map[string]interface{}) {
	// check if the room is valid
	room, err := h.roomService.GetRoom(c, roomID)
	if err != nil {
		conn.WriteJSON(gin.H{"type": "error", "message": "Room not found", "data": nil})
		return
	}

	// Get the game type from the parsed message
	gameTypeStr, ok := msg["game_type"].(string)
	if !ok {
		conn.WriteJSON(gin.H{"type": "error", "message": "Game type is required", "data": nil})
		return
	}

	gameType := model.GameType(gameTypeStr)

	// Handle the game choice through the service
	if err := h.roomService.HandleGameChosen(c, room, player, gameType); err != nil {
		conn.WriteJSON(gin.H{"type": "error", "message": "Failed to handle game choice", "data": nil})
		return
	}

	// Send confirmation to the player who chose the game
	conn.WriteJSON(gin.H{
		"type":    "game_chosen_confirmation",
		"message": "Your game choice has been recorded",
		"data": map[string]interface{}{
			"game_type": gameType,
		},
	})
}

// handleGameAccepted handles a player accepting a game
func (h *RoomHandler) handleGameAccepted(c *gin.Context, conn *websocket.Conn, roomID string, player model.Player, msg map[string]interface{}) {
	// check if the room is valid
	room, err := h.roomService.GetRoom(c, roomID)
	if err != nil {
		conn.WriteJSON(gin.H{"type": "error", "message": "Room not found", "data": nil})
		return
	}

	// Get the game type from the parsed message
	gameTypeStr, ok := msg["game_type"].(string)
	if !ok {
		conn.WriteJSON(gin.H{"type": "error", "message": "Game type is required", "data": nil})
		return
	}

	gameType := model.GameType(gameTypeStr)

	// Handle the game choice through the service
	if err := h.roomService.HandleGameAccepted(c, room, player, gameType); err != nil {
		conn.WriteJSON(gin.H{"type": "error", "message": "Failed to handle game accepted", "data": nil})
		return
	}

}

// handleGameRejected handles a player rejecting a game
func (h *RoomHandler) handleGameRejected(c *gin.Context, conn *websocket.Conn, roomID string, player model.Player, msg map[string]interface{}) {
	// check if the room is valid
	room, err := h.roomService.GetRoom(c, roomID)
	if err != nil {
		conn.WriteJSON(gin.H{"type": "error", "message": "Room not found", "data": nil})
		return
	}

	// Get the game type from the parsed message
	gameTypeStr, ok := msg["game_type"].(string)
	if !ok {
		conn.WriteJSON(gin.H{"type": "error", "message": "Game type is required", "data": nil})
		return
	}

	gameType := model.GameType(gameTypeStr)

	// Handle the game choice through the service
	if err := h.roomService.HandleGameRejected(c, room, player, gameType); err != nil {
		conn.WriteJSON(gin.H{"type": "error", "message": "Failed to handle game rejected", "data": nil})
		return
	}
}

// handleGameMove handles a player making a move in the game
func (h *RoomHandler) handleGameMove(c *gin.Context, conn *websocket.Conn, roomID string, player model.Player, msg map[string]interface{}) {
	room, err := h.roomService.GetRoom(c, roomID)
	if err != nil {
		conn.WriteJSON(gin.H{"type": "error", "message": "Room not found", "data": nil})
		return
	}

	_, exists := room.Players[player.User.ID]
	if !exists {
		conn.WriteJSON(gin.H{"type": "error", "message": "Player not found in room", "data": nil})
		return
	}

	move, ok := msg["move"]
	if !ok {
		conn.WriteJSON(gin.H{"type": "error", "message": "Move is required", "data": nil})
		return
	}

	room.Game.MakeMove(player.User.ID, move)

	// Get the game type from the parsed message
	conn.WriteJSON(gin.H{
		"type":    "move_received",
		"message": "Move received, processing...",
		"data":    nil,
	})
}
