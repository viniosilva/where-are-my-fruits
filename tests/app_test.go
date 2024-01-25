package controller

import (
	"bytes"
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

		// defers
		defer db.Exec("DELETE FROM buckets")

		// given
		r := api.ConfigGin(config.Api.Host, config.Api.Port, logger, factory.HealthController, factory.BucketController)

		createBucketReq := presenters.CreateBucketReq{
			Name:     "Testing",
			Capacity: 1,
		}

		// cases
		getHealth(t, r)
		postBucket(t, r, createBucketReq)
	})
}

func getHealth(t *testing.T, r *gin.Engine) {
	// given
	wantCode := http.StatusOK
	wantBody := presenters.HealthCheckRes{
		Status: presenters.HealthCheckStatusUp,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/healthcheck", nil)

	var got presenters.HealthCheckRes

	// when
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &got)

	// then
	assert.Equal(t, wantCode, w.Code)
	assert.Equal(t, wantBody, got)
}

func postBucket(t *testing.T, r *gin.Engine, data presenters.CreateBucketReq) presenters.BucketRes {
	// given
	wantCode := http.StatusCreated
	wantBody := presenters.BucketRes{
		Name:     "Testing",
		Capacity: 1,
	}

	body, _ := json.Marshal(data)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/buckets", bytes.NewReader(body))

	var got presenters.BucketRes

	// when
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &got)

	wantBody.ID = got.ID
	wantBody.CreatedAt = got.CreatedAt

	// then
	assert.Equal(t, wantCode, w.Code)
	assert.Equal(t, wantBody, got)

	return got
}

// Refers: https://gin-gonic.com/docs/testing
