package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kaviraj-j/duoplay/internal/middleware"
	"github.com/kaviraj-j/duoplay/internal/model"
	"github.com/kaviraj-j/duoplay/internal/service"
)

// WebSocket message structs
type WSMessage struct {
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

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
		conn.WriteJSON(WSMessage{Type: "error", Message: "Unauthorized", Data: nil})
		conn.Close()
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
		conn.WriteJSON(WSMessage{Type: "error", Message: err.Error(), Data: nil})
		return
	}

	// Send room created message via WebSocket
	conn.WriteJSON(WSMessage{
		Type:    "room_created",
		Message: "Room created successfully",
		Data:    room.GetRoomResponse(),
	})

	go h.handleWebSocketMessages(c, conn, room.ID, player)
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
		conn.WriteJSON(WSMessage{Type: "error", Message: "Room not found", Data: nil})
		conn.Close()
		return
	}

	if err := h.roomService.AddPlayer(c, roomID, player); err != nil {
		conn.Close()
		conn.WriteJSON(WSMessage{Type: "error", Message: err.Error(), Data: nil})
		return
	}

	// change room status to game selection
	room.Status = model.RoomStatusGameSelection
	err = h.roomService.UpdateRoom(c, *room)
	if err != nil {
		conn.WriteJSON(WSMessage{Type: "error", Message: "Failed to update room", Data: nil})
		conn.Close()
		return
	}

	// Send joined room message via WebSocket
	fmt.Printf("Sending joined_room message to player %s for room %s\n", player.User.ID, roomID)
	conn.WriteJSON(WSMessage{
		Type:    "joined_room",
		Message: "Joined room successfully",
		Data:    room.GetRoomResponse(),
	})

	// Send message to the other player
	var otherPlayer model.Player
	for userId, p := range room.Players {
		if userId != user.ID {
			otherPlayer = p
			break
		}
	}
	otherPlayer.Conn.WriteJSON(WSMessage{
		Type:    "joined_room",
		Message: "Opponent has joined room",
	})

	// Handle WebSocket connection
	go h.handleWebSocketMessages(c, conn, roomID, player)
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
			conn.WriteJSON(WSMessage{Type: "error", Message: "User is already in queue", Data: nil})
		} else {
			conn.WriteJSON(WSMessage{Type: "error", Message: "Failed to join queue", Data: nil})
		}
		return
	}

	// Send queue joined message via WebSocket
	conn.WriteJSON(WSMessage{
		Type:    "queue_joined",
		Message: "Successfully joined queue",
		Data:    nil,
	})

	// Handle queue WebSocket connection
	go h.handleQueueConnection(c, conn, user.ID)
}

// handleQueueConnection handles WebSocket communication for queue waiting
func (h *RoomHandler) handleQueueConnection(ctx *gin.Context, conn *websocket.Conn, userID string) {
	defer conn.Close()

	// Send initial message
	conn.WriteJSON(WSMessage{
		Type:    "queue_joined",
		Message: "Waiting for an opponent to connect...",
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

	c.JSON(http.StatusOK, gin.H{"type": "success", "message": "Room fetched successfully", "data": room.GetRoomResponse()})
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
	userInterface, exists := c.Get(middleware.AuthorizationPayloadKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"type": "error", "message": "Unauthorized", "data": nil})
		return
	}
	user := userInterface.(*model.User)
	// Get the other user's connection
	room, _ := h.roomService.GetRoom(c, roomID)
	oppositePlayer, err := h.roomService.GetOppositePlayer(c, room.Players, user.ID)

	if err := h.roomService.LeaveRoom(c, roomID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"type": "error", "message": "Failed to leave room", "data": nil})
		return
	}

	if err == nil {
		oppositePlayer.Conn.WriteJSON(WSMessage{
			Type:    "opponent_left",
			Message: "Opponent has left the room",
		})
	}

	c.JSON(http.StatusOK, gin.H{"type": "success", "message": "Left room", "data": nil})
}
