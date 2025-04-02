package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kaviraj-j/duoplay/internal/service"
)

type AuthMiddleWare struct {
	userService service.UserService
}

// auth related constants
const (
	authHeaderKey           = "Authorization"
	authType                = "Bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func NewAuthMiddleware(userService service.UserService) *AuthMiddleWare {
	return &AuthMiddleWare{
		userService: userService,
	}
}

func (authMiddleware *AuthMiddleWare) IsLoggedIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// get token from header
		authorizationHeader := ctx.GetHeader(authHeaderKey)
		if len(authorizationHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"type":    "error",
				"message": "authorization header not provided",
			})
			return
		}

		fields := strings.Split(authorizationHeader, " ")
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"type":    "error",
				"message": "invalid authorization header format",
			})
			return
		}

		if authType != fields[0] {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"type":    "error",
				"message": fmt.Sprintf("invalid auth type: %s, required: %s", fields[0], authType),
			})
			return
		}

		// validate token
		tokenStringFromHeader := fields[1]
		user, err := authMiddleware.userService.ValidateToken(ctx, tokenStringFromHeader)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"type":    "error",
				"message": err.Error(),
			})
			return
		}

		ctx.Set(AuthorizationPayloadKey, user)
		ctx.Next()
	}
}
