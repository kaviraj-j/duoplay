package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaviraj-j/duoplay/internal/model"
	"github.com/kaviraj-j/duoplay/internal/service"
)

type RoomMiddleWare struct {
	roomService *service.RoomService
}

func NewRoomMiddleware(roomService *service.RoomService) *RoomMiddleWare {
	return &RoomMiddleWare{
		roomService: roomService,
	}
}

func (roomMiddleware *RoomMiddleWare) IsRoomOwner() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roomID := ctx.Param("roomID")
		userPayload, _ := ctx.Get(AuthorizationPayloadKey)
		user := userPayload.(model.User)

		room, err := roomMiddleware.roomService.GetRoom(ctx, roomID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"type": "error", "message": "Room not found"})
			ctx.Abort()
		}

		_, exists := room.Players[user.ID]

		if !exists {
			ctx.JSON(http.StatusForbidden, gin.H{"type": "error", "message": "You are not a player in this room"})
			ctx.Abort()
		}

		ctx.Next()
	}
}
