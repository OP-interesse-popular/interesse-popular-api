package app

import (
	"github.com/OP-interesse-popular/interesse-popular-api/config"
	"github.com/gin-gonic/gin"
)

// App wires the HTTP server and the application's dependencies.
// Add domain handlers and services here as the API grows.
type App struct {
	Router *gin.Engine
	Config config.Config
}

func New(router *gin.Engine, cfg config.Config) *App {
	return &App{
		Router: router,
		Config: cfg,
	}
}
