package app

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kaviraj-j/duoplay/internal/config"
	"github.com/kaviraj-j/duoplay/internal/handler"
	"github.com/kaviraj-j/duoplay/internal/middleware"
	"github.com/kaviraj-j/duoplay/internal/repository"
	"github.com/kaviraj-j/duoplay/internal/service"
)

type App struct {
	config         config.Config
	userHandler    *handler.UserHandler
	roomHandler    *handler.RoomHandler
	gameHandler    *handler.GameHandler
	authMiddleware *middleware.AuthMiddleWare
	roomMiddleware *middleware.RoomMiddleWare
}

// creates new app
func Create(config config.Config) (*App, error) {
	// get user repo, service, and handler
	userRepository := repository.NewUserRepository()
	userService, err := service.CreateUserService(userRepository, []byte(config.JwtSecret))
	if err != nil {
		return nil, err
	}
	userHandler := handler.NewUserHandler(userService)
	authMiddleware := middleware.NewAuthMiddleware(userService)

	gameRepo := repository.NewGameRepository()
	gameService := service.NewGameService(gameRepo)
	gameHandler := handler.NewGameHandler(gameService)

	// room repo, service, handler and middleware
	roomRepo := repository.NewRoomRepository()
	queueRepo := repository.NewQueueRepository()
	roomService := service.NewRoomService(roomRepo, queueRepo, userRepository)
	roomHandler := handler.NewRoomHandler(roomService)
	roomMiddleware := middleware.NewRoomMiddleware(roomService)

	app := &App{
		config:         config,
		userHandler:    userHandler,
		authMiddleware: authMiddleware,
		roomHandler:    roomHandler,
		roomMiddleware: roomMiddleware,
		gameHandler:    gameHandler,
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

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// setup health route handler
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message":   "OK",
			"timestamp": time.Now(),
		})
	})

	// user related routes
	router.POST("/user", app.userHandler.NewUser)
	router.GET("/user/me", app.authMiddleware.IsAuthenticated(), app.userHandler.LoggedInUserDetails)

	// room routes
	router.GET("/room/join", app.authMiddleware.IsAuthenticated(), app.roomHandler.NewRoom)
	router.GET("/room/:roomID", app.authMiddleware.IsAuthenticated(), app.roomMiddleware.IsRoomOwner(), app.roomHandler.GetRoom)
	router.GET("/room/:roomID/join", app.authMiddleware.IsAuthenticated(), app.roomHandler.JoinRoom)
	router.GET("/room/joinQueue", app.authMiddleware.IsAuthenticated(), app.roomHandler.JoinWaitingQueue)

	// game routes
	router.GET("/game/list", app.gameHandler.GetGamesList)

}
