package app

import (
	"net/http"

	"github.com/TiB-Software/go-commons/commonsmodels"
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
		commonsutils.SendResponse(c, commonsmodels.NewResponseSuccess(), http.StatusOK)
	})
	api.GET("/parliamentarians/right-wing", a.Parliamentarians.ListRightWing())
	api.GET("/parliamentarians/right-wing/camara", a.Parliamentarians.ListRightWingCamara())
	api.GET("/parliamentarians/right-wing/senado", a.Parliamentarians.ListRightWingSenado())
}
