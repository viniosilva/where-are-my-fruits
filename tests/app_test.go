package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

		defer func() {
			db.SQL.Exec("DELETE FROM fruits")
			db.SQL.Exec("DELETE FROM buckets")
		}()

		// given
		r := api.ConfigGin(config.Api.Host, config.Api.Port, logger, factory.HealthController, factory.BucketController, factory.FruitController)

		// cases
		getHealth(t, r)

		// case: create bucket
		bucket := createBucket(t, r, presenters.CreateBucketReq{
			Name:     "Medium fruits",
			Capacity: 2,
		}, http.StatusCreated,
			&presenters.BucketRes{
				Name:     "Medium fruits",
				Capacity: 2,
			})

		// case: create apple fruit expires in 1h out of the bucket
		fruit := createFruit(t, r, presenters.CreateFruitReq{
			Name:      "Apple",
			Price:     decimal.NewFromFloat32(1.99),
			ExpiresIn: "1h",
		}, http.StatusCreated,
			&presenters.FruitRes{
				Name:  "Apple",
				Price: decimal.NewFromFloat32(1.99),
			})

		// case: add apple to the bucket
		addFruitOnBucket(t, r, fruit.ID, bucket.ID, http.StatusOK)

		// case: create melon fruit expires in 1s inside the bucket
		createFruit(t, r, presenters.CreateFruitReq{
			Name:      "Melon",
			Price:     decimal.NewFromFloat32(3.50),
			BucketID:  &bucket.ID,
			ExpiresIn: "1s",
		}, http.StatusCreated,
			&presenters.FruitRes{
				Name:     "Melon",
				Price:    decimal.NewFromFloat32(3.50),
				BucketID: &bucket.ID,
			})

		// case: list buckets with all fruits
		listBuckets(t, r, http.StatusOK, &presenters.BucketsFruitsRes{
			Data: []presenters.BucketFruitsRes{
				{
					ID:          bucket.ID,
					Name:        "Medium fruits",
					Capacity:    2,
					TotalFruits: 2,
					TotalPrice:  decimal.NewFromFloat32(5.49),
					Percent:     "100.00%",
				},
			},
		})

		// wait for melon expiration
		time.Sleep(1 * time.Second)

		// case: list buckets with one fruit
		listBuckets(t, r, http.StatusOK, &presenters.BucketsFruitsRes{
			Data: []presenters.BucketFruitsRes{
				{
					ID:          bucket.ID,
					Name:        "Medium fruits",
					Capacity:    2,
					TotalFruits: 1,
					TotalPrice:  decimal.NewFromFloat32(1.99),
					Percent:     "50.00%",
				},
			},
		})

		// case: try remove bucket, but fail for it be full
		deleteBucket(t, r, bucket.ID, http.StatusBadRequest)

		// case: create abacato fruit expires in 1s inside the bucket
		createFruit(t, r, presenters.CreateFruitReq{
			Name:      "Abacato",
			Price:     decimal.NewFromFloat32(7.50),
			BucketID:  &bucket.ID,
			ExpiresIn: "1s",
		}, http.StatusCreated, &presenters.FruitRes{
			Name:     "Abacato",
			Price:    decimal.NewFromFloat32(7.50),
			BucketID: &bucket.ID,
		})

		// case: try add another abacato to bucket, but fail for it be full
		createFruit(t, r, presenters.CreateFruitReq{
			Name:      "Abacato",
			Price:     decimal.NewFromFloat32(7.50),
			BucketID:  &bucket.ID,
			ExpiresIn: "1s",
		}, http.StatusBadRequest, nil)

		// case: remove apple from bucket
		removeFruitFromBucket(t, r, fruit.ID, bucket.ID, http.StatusOK)

		// wait for abacato expiration
		time.Sleep(1 * time.Second)

		// case: list empty buckets
		listBuckets(t, r, http.StatusOK, &presenters.BucketsFruitsRes{
			Data: []presenters.BucketFruitsRes{
				{
					ID:          bucket.ID,
					Name:        "Medium fruits",
					Capacity:    2,
					TotalFruits: 0,
					TotalPrice:  decimal.NewFromInt32(0),
					Percent:     "0.00%",
				},
			},
		})

		// case: delete apple from bucket
		deleteFruit(t, r, fruit.ID, http.StatusOK)

		// case: delete bucket
		deleteBucket(t, r, bucket.ID, http.StatusOK)
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

func createBucket(t *testing.T, r *gin.Engine, data presenters.CreateBucketReq, wantCode int, wantBody *presenters.BucketRes) presenters.BucketRes {
	// given
	body, _ := json.Marshal(data)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/buckets", bytes.NewReader(body))

	var got presenters.BucketRes

	// when
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &got)

	// then
	assert.Equal(t, wantCode, w.Code)

	if wantBody != nil {
		wantBody.ID = got.ID
		wantBody.CreatedAt = got.CreatedAt

		assert.Equal(t, *wantBody, got)
	}

	return got
}

func createFruit(t *testing.T, r *gin.Engine, data presenters.CreateFruitReq, wantCode int, wantBody *presenters.FruitRes) presenters.FruitRes {
	// given
	body, _ := json.Marshal(data)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fruits", bytes.NewReader(body))

	var got presenters.FruitRes

	// when
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &got)

	// then
	assert.Equal(t, wantCode, w.Code)

	if wantBody != nil {
		wantBody.ID = got.ID
		wantBody.CreatedAt = got.CreatedAt
		wantBody.ExpiresAt = got.ExpiresAt
		wantBody.BucketID = got.BucketID

		assert.Equal(t, *wantBody, got)
	}

	return got
}

func addFruitOnBucket(t *testing.T, r *gin.Engine, fruitID, bucketID int64, wantCode int) {
	// given
	path := fmt.Sprintf("/api/v1/fruits/%d/buckets/%d", fruitID, bucketID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", path, nil)

	// when
	r.ServeHTTP(w, req)

	// then
	assert.Equal(t, wantCode, w.Code)
}

func removeFruitFromBucket(t *testing.T, r *gin.Engine, fruitID, bucketID int64, wantCode int) {
	// given
	path := fmt.Sprintf("/api/v1/fruits/%d/buckets/%d", fruitID, bucketID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", path, nil)

	// when
	r.ServeHTTP(w, req)

	// then
	assert.Equal(t, wantCode, w.Code)
}

func deleteFruit(t *testing.T, r *gin.Engine, fruitID int64, wantCode int) {
	// given
	path := fmt.Sprintf("/api/v1/fruits/%d", fruitID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", path, nil)

	// when
	r.ServeHTTP(w, req)

	// then
	assert.Equal(t, wantCode, w.Code)
}

func deleteBucket(t *testing.T, r *gin.Engine, bucketID int64, wantCode int) {
	// given
	path := fmt.Sprintf("/api/v1/buckets/%d", bucketID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", path, nil)

	// when
	r.ServeHTTP(w, req)

	// then
	assert.Equal(t, wantCode, w.Code)
}

func listBuckets(t *testing.T, r *gin.Engine, wantCode int, wantBody *presenters.BucketsFruitsRes) {
	// given
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/buckets", nil)

	var got presenters.BucketsFruitsRes

	// when
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &got)

	// then
	assert.Equal(t, wantCode, w.Code)

	if wantBody != nil {
		assert.Equal(t, *wantBody, got)
	}
}

// Refers: https://gin-gonic.com/docs/testing
