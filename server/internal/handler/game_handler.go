package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kaviraj-j/duoplay/internal/service"
)

type GameHandler struct {
	gameService *service.GameService
	upgrader    websocket.Upgrader
}

func NewGameHandler(s *service.GameService) *GameHandler {
	return &GameHandler{
		gameService: s,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // TODO: need to implement a proper origin checking
			},
		},
	}
}

func (h *GameHandler) GetGamesList(c *gin.Context) {
	gamesList, err := h.gameService.GetGamesList(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, gin.H{"type": "success", "data": gamesList})
}
