package app

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kaviraj-j/duoplay/internal/config"
	"github.com/kaviraj-j/duoplay/internal/handler"
	"github.com/kaviraj-j/duoplay/internal/repository"
	"github.com/kaviraj-j/duoplay/internal/service"
)

type App struct {
	config      *config.Config
	userHandler *handler.UserHandler
}

// creates new app
func Create(config *config.Config) (*App, error) {
	// get user repo, service, and handler
	userRepository := repository.NewUserRepository()
	userService, err := service.CreateUserService(userRepository, []byte(config.JwtSecret))
	if err != nil {
		return nil, err
	}
	userHandler := handler.NewUserHandler(&userService)
	app := &App{
		config:      config,
		userHandler: userHandler,
	}
	return app, nil
}

// Run will setup routes and run the app
func (app *App) Run() {
	router := gin.Default()
	app.setupRouter(router)

	router.Run(app.config.ServerAddress)
}

func (app *App) setupRouter(router *gin.Engine) {

	// setup health route handler
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message":   "OK",
			"timestamp": time.Now(),
		})
	})

	// user related routes
	router.POST("/user", app.userHandler.NewUser)

}
