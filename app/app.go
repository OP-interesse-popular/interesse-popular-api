package app

import (
	"github.com/OP-interesse-popular/interesse-popular-api/cache"
	"github.com/OP-interesse-popular/interesse-popular-api/config"
	"github.com/OP-interesse-popular/interesse-popular-api/handlers"
	"github.com/OP-interesse-popular/interesse-popular-api/services"
	"github.com/gin-gonic/gin"
)

// App wires the HTTP server and the application's dependencies.
// Add domain handlers and services here as the API grows.
type App struct {
	Router           *gin.Engine
	Config           config.Config
	Parliamentarians *handlers.Parliamentarians
}

func New(router *gin.Engine, cfg config.Config) *App {
	parliamentarianCache, err := cache.NewRedisParliamentarianCache(cfg.RedisURL)
	if err != nil {
		panic(err)
	}
	parliamentariansService := services.NewParliamentarians(parliamentarianCache)

	return &App{
		Router:           router,
		Config:           cfg,
		Parliamentarians: handlers.NewParliamentarians(parliamentariansService),
	}
}
