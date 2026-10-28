package app

import (
	"net/http"

	commonsutils "github.com/TiB-Software/go-commons/utils"
	"github.com/gin-gonic/gin"
)

func (a *App) RegisterRoutes() {
	a.Router.Use(
		commonsutils.CORS(commonsutils.CORSConfig{
			AllowedOrigins:   a.Config.CORSAllowedOrigins,
			AllowCredentials: true,
		}),
		commonsutils.SecurityHeaders(commonsutils.SecurityHeadersConfig{
			Production: a.Config.IsProduction(),
		}),
		commonsutils.RequireHTTPS(commonsutils.HTTPSConfig{
			Production:     a.Config.IsProduction(),
			TrustedProxies: a.Config.TrustedProxies,
		}),
	)

	api := a.Router.Group("/api/v1")
	api.GET("/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
}
