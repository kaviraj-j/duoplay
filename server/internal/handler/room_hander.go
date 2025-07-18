package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kaviraj-j/duoplay/internal/middleware"
	"github.com/kaviraj-j/duoplay/internal/model"
	"github.com/kaviraj-j/duoplay/internal/service"
)

type RoomHandler struct {
	roomService *service.RoomService
	upgrader    websocket.Upgrader
}

func NewRoomHandler(s *service.RoomService) *RoomHandler {
	return &RoomHandler{
		roomService: s,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // TODO: need to implement a proper origin checking
			},
		},
	}
}

// NewRoom creates a new game room
func (h *RoomHandler) NewRoom(c *gin.Context) {
	room, err := h.roomService.CreateRoom(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":    "error",
			"message": "Failed to create room",
			"data":    nil,
		})
		return
	}

	// upgrade http connection to websocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"type": "error", "message": "Could not upgrade connection", "data": nil})
		return
	}

	userInterface, exists := c.Get(middleware.AuthorizationPayloadKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"type": "error", "message": "Unauthorized", "data": nil})
		return
	}
	user := userInterface.(*model.User)

	// create player and add to room
	player := model.Player{
		User: *user,
		Conn: conn,
	}

	if err := h.roomService.AddPlayer(c, room.ID, player); err != nil {
		conn.Close()
		conn.WriteJSON(gin.H{"type": "error", "message": err.Error(), "data": nil})
		return
	}

	go h.handleWebSocketMessages(c, conn, room.ID, player)

	c.JSON(http.StatusCreated, gin.H{
		"type":    "success",
		"message": "Room created successfully",
		"data":    room,
	})
}

// JoinRoom handles player joining a room via WebSocket
func (h *RoomHandler) JoinRoom(c *gin.Context) {
	roomID := c.Param("roomID")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"type": "error", "message": "Room ID is required", "data": nil})
		return
	}

	// Get user from auth middleware
	userInterface, exists := c.Get(middleware.AuthorizationPayloadKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"type": "error", "message": "Unauthorized", "data": nil})
		return
	}
	user := userInterface.(*model.User)

	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"type": "error", "message": "Could not upgrade connection", "data": nil})
		return
	}

	// Create player and add to room
	player := model.Player{
		User: *user,
		Conn: conn,
	}

	room, err := h.roomService.GetRoom(c, roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"type": "error", "message": "Room not found", "data": nil})
		return
	}

	if err := h.roomService.AddPlayer(c, roomID, player); err != nil {
		conn.Close()
		conn.WriteJSON(gin.H{"type": "error", "message": err.Error(), "data": nil})
		return
	}

	// change room status to game selection
	room.Status = model.RoomStatusGameSelection
	err = h.roomService.UpdateRoom(c, *room)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"type": "error", "message": "Failed to update room", "data": nil})
		return
	}

	// Handle WebSocket connection
	go h.handleWebSocketMessages(c, conn, roomID, player)

	c.JSON(http.StatusOK, gin.H{"type": "success", "message": "Joined room", "data": room})
}

func (h *RoomHandler) JoinWaitingQueue(c *gin.Context) {
	userInterface, exists := c.Get(middleware.AuthorizationPayloadKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"type": "error", "message": "Unauthorized", "data": nil})
		return
	}
	user := userInterface.(*model.User)

	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"type": "error", "message": "Could not upgrade connection", "data": nil})
		return
	}

	// Add user to queue with WebSocket connection
	if err := h.roomService.JoinQueue(c, user.ID, conn); err != nil {
		conn.Close()
		if err == service.ErrorUserAlreadyInQueue {
			c.JSON(http.StatusBadRequest, gin.H{"type": "error", "message": "User is already in queue", "data": nil})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"type": "error", "message": "Failed to join queue", "data": nil})
		}
		return
	}

	// Handle queue WebSocket connection
	go h.handleQueueConnection(c, conn, user.ID)
}

// handleQueueConnection handles WebSocket communication for queue waiting
func (h *RoomHandler) handleQueueConnection(ctx *gin.Context, conn *websocket.Conn, userID string) {
	defer conn.Close()

	// Send initial message
	conn.WriteJSON(gin.H{
		"type":    "queue_joined",
		"message": "Waiting for an opponent to connect...",
	})

	for {
		// Read message from WebSocket
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			// Handle disconnection - remove from queue
			h.roomService.RemoveFromQueue(ctx, userID)
			return
		}

		// Handle different message types
		switch messageType {
		case websocket.TextMessage:
			// Handle text messages if needed
			continue
		case websocket.CloseMessage:
			// Remove from queue on close
			h.roomService.RemoveFromQueue(ctx, userID)
			return
		}
	}
}

// GetRoom returns room details
func (h *RoomHandler) GetRoom(c *gin.Context) {
	roomID := c.Param("roomID")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"type": "error", "message": "Room ID is required", "data": nil})
		return
	}

	room, err := h.roomService.GetRoom(c, roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"type": "error", "message": "Room not found", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"type": "success", "message": "Room fetched successfully", "data": room})
}

// StartGame initiates the game in the room
func (h *RoomHandler) StartGame(c *gin.Context) {
	roomID := c.Param("roomID")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"type": "error", "message": "Room ID is required", "data": nil})
		return
	}

	if err := h.roomService.StartGame(c, roomID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"type": "error", "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"type": "success", "message": "Game started", "data": nil})
}

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
		case model.MessageTypeGameChosen:
			h.handleGameChosen(c, conn, roomID, player, msg)
		case model.MessageTypeGameAccepted:
			h.handleGameAccepted(c, conn, roomID, player, msg)
		case model.MessageTypeGameMove, model.MessageTypeMoveMade:
			h.handleGameMove(c, conn, roomID, player, msg)
		case model.MessageTypeRejectGame:
		default:
			conn.WriteJSON(gin.H{"type": "error", "message": "Unknown message type", "data": nil})
		}
	}
}

func (h *RoomHandler) LeaveWaitingQueue(c *gin.Context) {
	userInterface, exists := c.Get(middleware.AuthorizationPayloadKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"type": "error", "message": "Unauthorized", "data": nil})
		return
	}
	user := userInterface.(*model.User)

	if err := h.roomService.RemoveFromQueue(c, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"type": "error", "message": "Failed to leave queue", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"type": "success", "message": "Left queue", "data": nil})
}

func (h *RoomHandler) LeaveRoom(c *gin.Context) {
	roomID := c.Param("roomID")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"type": "error", "message": "Room ID is required", "data": nil})
		return
	}

	if err := h.roomService.LeaveRoom(c, roomID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"type": "error", "message": "Failed to leave room", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"type": "success", "message": "Left room", "data": nil})
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

// handleGameMove handles a player making a move in the game
func (h *RoomHandler) handleGameMove(c *gin.Context, conn *websocket.Conn, roomID string, player model.Player, msg map[string]interface{}) {
	// This will be implemented when the actual game logic is added
	// For now, just acknowledge the move
	conn.WriteJSON(gin.H{
		"type":    "move_received",
		"message": "Move received, processing...",
		"data":    nil,
	})
}
