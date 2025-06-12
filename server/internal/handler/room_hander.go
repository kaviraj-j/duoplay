package handler

import (
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
			"error": "Failed to create room",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"room_id": room.ID,
		"message": "Room created successfully",
	})
}

// JoinRoom handles player joining a room via WebSocket
func (h *RoomHandler) JoinRoom(c *gin.Context) {
	roomID := c.Param("roomID")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	// Get user from auth middleware
	userInterface, exists := c.Get(middleware.AuthorizationPayloadKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(*model.User)

	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not upgrade connection"})
		return
	}

	// Create player and add to room
	player := model.Player{
		User: *user,
		Conn: conn,
	}

	if err := h.roomService.AddPlayer(c, roomID, player); err != nil {
		conn.Close()
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	// Handle WebSocket connection
	go h.handleGameConnection(c, conn, roomID, player)
}

// GetRoom returns room details
func (h *RoomHandler) GetRoom(c *gin.Context) {
	roomID := c.Param("roomID")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	room, err := h.roomService.GetRoom(c, roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	c.JSON(http.StatusOK, room)
}

// StartGame initiates the game in the room
func (h *RoomHandler) StartGame(c *gin.Context) {
	roomID := c.Param("roomID")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	if err := h.roomService.StartGame(c, roomID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game started"})
}

// handleGameConnection handles WebSocket communication for a game session
func (h *RoomHandler) handleGameConnection(ctx *gin.Context, conn *websocket.Conn, roomID string, player model.Player) {
	defer conn.Close()

	for {
		// Read message from WebSocket
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			// Handle disconnection
			// h.roomService.HandlePlayerDisconnect(ctx, roomID, player.User.ID)
			return
		}

		// Handle different message types
		switch messageType {
		case websocket.TextMessage:

			continue
		case websocket.CloseMessage:
			return
		}
	}
}
