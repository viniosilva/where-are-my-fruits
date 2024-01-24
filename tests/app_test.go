package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/viniosilva/where-are-my-fruits/api"
	"github.com/viniosilva/where-are-my-fruits/internal/controllers/presenters"
	"github.com/viniosilva/where-are-my-fruits/internal/factories"
	"github.com/viniosilva/where-are-my-fruits/internal/infra"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/api/healthcheck", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}

func TestApp(t *testing.T) {
	t.Run("test app", func(t *testing.T) {
		// setup
		logger := infra.ConfigLogger()

		config, err := infra.GetConfig("..")
		require.Nil(t, err)

		db, err := infra.ConfigDB(config.MySQL.Username, config.MySQL.Password, config.MySQL.Host,
			config.MySQL.Port, config.MySQL.Database, config.MySQL.ConnMaxLifetime,
			config.MySQL.MaxIdleConns, config.MySQL.MaxOpenConns)
		require.Nil(t, err)

		factory, err := factories.Build(db, logger)
		require.Nil(t, err)

		// given
		r := api.ConfigGin(config.Api.Host, config.Api.Port, logger, factory.HealthController)

		// cases
		getHealth(t, r)
	})
}

func getHealth(t *testing.T, r *gin.Engine) {
	// given
	wantCode := http.StatusOK
	wantBody := presenters.HealthCheckResponse{
		Status: presenters.HealthCheckStatusUp,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/healthcheck", nil)

	var got presenters.HealthCheckResponse

	// when
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &got)

	// then
	assert.Equal(t, wantCode, w.Code)
	assert.Equal(t, wantBody, got)
}

// Refers: https://gin-gonic.com/docs/testing
