package app

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kaviraj-j/duoplay/internal/config"
)

type App struct {
	config *config.Config
}

// creates new app
func Create(config *config.Config) (*App, error) {
	app := &App{
		config: config,
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

}
