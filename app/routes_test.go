package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OP-interesse-popular/interesse-popular-api/config"
	"github.com/gin-gonic/gin"
)

func TestHealthcheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	application := New(router, config.Config{
		CORSAllowedOrigins: []string{"http://localhost:3000"},
		TrustedProxies:     []string{"127.0.0.1", "::1"},
	})
	application.RegisterRoutes()

	request := httptest.NewRequest(http.MethodGet, "/api/v1/healthcheck", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("healthcheck status = %d, want %d", response.Code, http.StatusOK)
	}
}
