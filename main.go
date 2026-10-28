package main

import (
	"fmt"
	"os"

	"github.com/OP-interesse-popular/interesse-popular-api/app"
	"github.com/OP-interesse-popular/interesse-popular-api/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	cfg := config.Load()
	router := gin.Default()
	if err := router.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		panic(err)
	}

	application := app.New(router, cfg)
	application.RegisterRoutes()

	fmt.Printf("Starting server on port :%s\n", cfg.Port)
	if err := application.Router.Run(":" + cfg.Port); err != nil {
		panic(err)
	}
}
