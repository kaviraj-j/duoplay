package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaviraj-j/duoplay/internal/middleware"
	"github.com/kaviraj-j/duoplay/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: *userService,
	}
}

func (handler *UserHandler) NewUser(ctx *gin.Context) {
	var newUserRequestDetails struct {
		Name string `json:"name" binding:"required"`
	}
	if err := ctx.ShouldBindBodyWithJSON(&newUserRequestDetails); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "error while parsing data",
		})
		return
	}

	user, token, err := handler.userService.RegisterUser(ctx, newUserRequestDetails.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error while creating user",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user":    user,
		"token":   token,
		"message": "new user created successfully",
	})
}

func (handler *UserHandler) LoggedInUserDetails(ctx *gin.Context) {
	user, ok := ctx.Get(middleware.AuthorizationPayloadKey)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"type":    "error",
			"message": "user is not authorized",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"type": "success",
		"data": user,
	})
}
