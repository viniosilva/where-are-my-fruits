package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
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

		validate := infra.NewValidator()

		factory, err := factories.Build(db, logger, validate)
		require.Nil(t, err)

		// defers
		defer db.SQL.Exec("DELETE FROM fruits")
		defer db.SQL.Exec("DELETE FROM buckets")

		// given
		r := api.ConfigGin(config.Api.Host, config.Api.Port, logger, factory.HealthController, factory.BucketController, factory.FruitController)

		createBucketReq := presenters.CreateBucketReq{
			Name:     "Testing",
			Capacity: 2,
		}

		price, _ := decimal.NewFromString("1.99")
		createFruitReq := presenters.CreateFruitReq{
			Name:      "Testing",
			Price:     price,
			ExpiresIn: "1h",
		}

		// cases
		getHealth(t, r)

		bucket := createBucket(t, r, createBucketReq)
		fruit := createFruit(t, r, createFruitReq)

		addFruitOnBucket(t, r, fruit.ID, bucket.ID)

		createFruitReq.BucketID = &bucket.ID
		createFruit(t, r, createFruitReq)

		removeFruitFromBucket(t, r, fruit.ID, bucket.ID)
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

func createBucket(t *testing.T, r *gin.Engine, data presenters.CreateBucketReq) presenters.BucketRes {
	// given
	wantCode := http.StatusCreated
	wantBody := presenters.BucketRes{
		Name:     "Testing",
		Capacity: 2,
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

func createFruit(t *testing.T, r *gin.Engine, data presenters.CreateFruitReq) presenters.FruitRes {
	// given
	price, _ := decimal.NewFromString("1.99")
	wantCode := http.StatusCreated
	wantBody := presenters.FruitRes{
		Name:  "Testing",
		Price: price,
	}

	body, _ := json.Marshal(data)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fruits", bytes.NewReader(body))

	var got presenters.FruitRes

	// when
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &got)

	wantBody.ID = got.ID
	wantBody.CreatedAt = got.CreatedAt
	wantBody.ExpiresAt = got.ExpiresAt
	wantBody.BucketID = got.BucketID

	// then
	assert.Equal(t, wantCode, w.Code)
	assert.Equal(t, wantBody, got)

	return got
}

func addFruitOnBucket(t *testing.T, r *gin.Engine, fruitID, bucketID int64) {
	// given
	wantCode := http.StatusOK

	path := fmt.Sprintf("/api/v1/fruits/%d/buckets/%d", fruitID, bucketID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", path, nil)

	// when
	r.ServeHTTP(w, req)

	// then
	assert.Equal(t, wantCode, w.Code)
}

func removeFruitFromBucket(t *testing.T, r *gin.Engine, fruitID, bucketID int64) {
	// given
	wantCode := http.StatusOK

	path := fmt.Sprintf("/api/v1/fruits/%d/buckets/%d", fruitID, bucketID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", path, nil)

	// when
	r.ServeHTTP(w, req)

	// then
	assert.Equal(t, wantCode, w.Code)
}

// Refers: https://gin-gonic.com/docs/testing
