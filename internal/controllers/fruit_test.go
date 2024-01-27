package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/where-are-my-fruits/internal/controllers/presenters"
	"github.com/viniosilva/where-are-my-fruits/internal/dtos"
	"github.com/viniosilva/where-are-my-fruits/internal/exceptions"
	"github.com/viniosilva/where-are-my-fruits/internal/models"
	"github.com/viniosilva/where-are-my-fruits/mocks"
)

func TestFruitController_Create(t *testing.T) {
	now := time.Date(2000, 12, 31, 23, 59, 59, 0, time.Local)
	price, _ := decimal.NewFromString("1.99")
	expiresIn, _ := time.ParseDuration("1m")
	bucketID := int64(1)

	tests := map[string]struct {
		mock        func(service *mocks.MockFruitService)
		body        presenters.CreateFruitReq
		wantCode    int
		wantBody    presenters.FruitRes
		wantBodyErr presenters.ErrorRes
	}{
		"should be success": {
			mock: func(service *mocks.MockFruitService) {
				data := dtos.CreateFruitDto{
					Name:      "Testing",
					Price:     price,
					ExpiresIn: &expiresIn,
				}
				service.EXPECT().Create(gomock.Any(), data).Return(&models.Fruit{
					ID:        1,
					CreatedAt: now,
					Name:      "Testing",
					Price:     price,
					ExpiresAt: now.Add(expiresIn),
				}, nil)
			},
			body: presenters.CreateFruitReq{
				Name:      "Testing",
				Price:     price,
				ExpiresIn: "1m",
			},
			wantCode: http.StatusCreated,
			wantBody: presenters.FruitRes{
				ID:        1,
				CreatedAt: "2000-12-31 23:59:59",
				Name:      "Testing",
				Price:     price,
				ExpiresAt: "2001-01-01 00:00:59",
			},
		},
		"should be success when bucketID is setted": {
			mock: func(service *mocks.MockFruitService) {
				data := dtos.CreateFruitDto{
					Name:      "Testing",
					Price:     price,
					ExpiresIn: &expiresIn,
					BucketID:  &bucketID,
				}
				service.EXPECT().Create(gomock.Any(), data).Return(&models.Fruit{
					ID:        1,
					CreatedAt: now,
					Name:      "Testing",
					Price:     price,
					ExpiresAt: now.Add(expiresIn),
					BucketID:  &bucketID,
				}, nil)
			},
			body: presenters.CreateFruitReq{
				Name:      "Testing",
				Price:     price,
				ExpiresIn: "1m",
				BucketID:  &bucketID,
			},
			wantCode: http.StatusCreated,
			wantBody: presenters.FruitRes{
				ID:        1,
				CreatedAt: "2000-12-31 23:59:59",
				Name:      "Testing",
				Price:     price,
				ExpiresAt: "2001-01-01 00:00:59",
				BucketID:  &bucketID,
			},
		},
		"should throw validation exception": {
			mock: func(service *mocks.MockFruitService) {
				service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, exceptions.NewValidationException(validator.ValidationErrors{
					&mocks.FieldError{Itag: "error 1", Ins: "error 1"},
					&mocks.FieldError{Itag: "error 2", Ins: "error 2"},
				}))
			},
			body:     presenters.CreateFruitReq{},
			wantCode: http.StatusBadRequest,
			wantBodyErr: presenters.ErrorRes{
				Error: exceptions.ValidationExceptionName,
				Messages: []string{
					"Key: 'error 1' Error:Field validation for '' failed on the 'error 1' tag",
					"Key: 'error 2' Error:Field validation for '' failed on the 'error 2' tag",
				},
			},
		},
		"should throw foreign not found exception": {
			mock: func(service *mocks.MockFruitService) {
				service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, exceptions.NewForeignNotFoundException("Bucket not found"))
			},
			body: presenters.CreateFruitReq{
				Name:      "Testing",
				Price:     price,
				ExpiresIn: "1m",
				BucketID:  &bucketID,
			},
			wantCode: http.StatusBadRequest,
			wantBodyErr: presenters.ErrorRes{
				Error:   exceptions.ValidationExceptionName,
				Message: "Bucket not found",
			},
		},
		"should throw forbidden exception": {
			mock: func(service *mocks.MockFruitService) {
				service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, exceptions.NewForbiddenException("Bucket is full"))
			},
			body: presenters.CreateFruitReq{
				Name:      "Testing",
				Price:     price,
				ExpiresIn: "1m",
				BucketID:  &bucketID,
			},
			wantCode: http.StatusBadRequest,
			wantBodyErr: presenters.ErrorRes{
				Error:   exceptions.ValidationExceptionName,
				Message: "Bucket is full",
			},
		},
		"should throw internal server error": {
			mock: func(service *mocks.MockFruitService) {
				service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			body:        presenters.CreateFruitReq{},
			wantCode:    http.StatusInternalServerError,
			wantBodyErr: presenters.ErrorRes{Error: http.StatusText(http.StatusInternalServerError)},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			serviceMock := mocks.NewMockFruitService(ctrl)
			tt.mock(serviceMock)

			r := gin.Default()
			controller := NewFruit(serviceMock)

			path := "/api/v1/fruits"
			r.POST(path, controller.Create)

			var got presenters.FruitRes
			var gotErr presenters.ErrorRes

			// given
			body, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", path, bytes.NewReader(body))

			// when
			r.ServeHTTP(w, req)

			errBodyErr := json.Unmarshal(w.Body.Bytes(), &gotErr)
			json.Unmarshal(w.Body.Bytes(), &got)

			// then
			assert.Equal(t, tt.wantCode, w.Code)

			if errBodyErr != nil {
				assert.Equal(t, tt.wantBodyErr, gotErr)
				return
			}
			assert.Equal(t, tt.wantBody, got)
		})
	}
}

func TestFruitController_AddOnBucket(t *testing.T) {
	tests := map[string]struct {
		mock          func(service *mocks.MockFruitService)
		fruitIDParam  string
		bucketIDParam string
		wantCode      int
		wantBodyErr   presenters.ErrorRes
	}{
		"should be success": {
			mock: func(service *mocks.MockFruitService) {
				service.EXPECT().AddOnBucket(gomock.Any(), int64(1), int64(1)).Return(nil)
			},
			fruitIDParam:  "1",
			bucketIDParam: "1",
			wantCode:      http.StatusOK,
		},
		"should throw validation exception when fruitID is invalid": {
			mock:          func(service *mocks.MockFruitService) {},
			fruitIDParam:  "invalid",
			bucketIDParam: "1",
			wantCode:      http.StatusBadRequest,
			wantBodyErr: presenters.ErrorRes{
				Error:   exceptions.ValidationExceptionName,
				Message: "invalid fruitID",
			},
		},
		"should throw validation exception when bucketID is invalid": {
			mock:          func(service *mocks.MockFruitService) {},
			fruitIDParam:  "1",
			bucketIDParam: "invalid",
			wantCode:      http.StatusBadRequest,
			wantBodyErr: presenters.ErrorRes{
				Error:   exceptions.ValidationExceptionName,
				Message: "invalid bucketID",
			},
		},
		"should throw foreign not found exception": {
			mock: func(service *mocks.MockFruitService) {
				service.EXPECT().AddOnBucket(gomock.Any(), gomock.Any(), gomock.Any()).Return(exceptions.NewForeignNotFoundException("Bucket not found"))
			},
			fruitIDParam:  "1",
			bucketIDParam: "1",
			wantCode:      http.StatusBadRequest,
			wantBodyErr: presenters.ErrorRes{
				Error:   exceptions.ValidationExceptionName,
				Message: "Bucket not found",
			},
		},
		"should throw forbidden exception": {
			mock: func(service *mocks.MockFruitService) {
				service.EXPECT().AddOnBucket(gomock.Any(), gomock.Any(), gomock.Any()).Return(exceptions.NewForbiddenException("Bucket is full"))
			},
			fruitIDParam:  "1",
			bucketIDParam: "1",
			wantCode:      http.StatusBadRequest,
			wantBodyErr: presenters.ErrorRes{
				Error:   exceptions.ValidationExceptionName,
				Message: "Bucket is full",
			},
		},
		"should throw not found exception": {
			mock: func(service *mocks.MockFruitService) {
				service.EXPECT().AddOnBucket(gomock.Any(), gomock.Any(), gomock.Any()).Return(exceptions.NewNotFoundException("Fruit not found"))
			},
			fruitIDParam:  "1",
			bucketIDParam: "1",
			wantCode:      http.StatusBadRequest,
			wantBodyErr: presenters.ErrorRes{
				Error:   exceptions.NotFoundExceptionName,
				Message: "Fruit not found",
			},
		},
		"should throw internal server error": {
			mock: func(service *mocks.MockFruitService) {
				service.EXPECT().AddOnBucket(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
			},
			fruitIDParam:  "1",
			bucketIDParam: "1",
			wantCode:      http.StatusInternalServerError,
			wantBodyErr:   presenters.ErrorRes{Error: http.StatusText(http.StatusInternalServerError)},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			serviceMock := mocks.NewMockFruitService(ctrl)
			tt.mock(serviceMock)

			r := gin.Default()
			controller := NewFruit(serviceMock)

			r.POST("/api/v1/fruits/:fruitID/buckets/:bucketID", controller.AddOnBucket)

			var gotErr presenters.ErrorRes

			// given
			path := fmt.Sprintf("/api/v1/fruits/%s/buckets/%s", tt.fruitIDParam, tt.bucketIDParam)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", path, nil)

			// when
			r.ServeHTTP(w, req)

			errBodyErr := json.Unmarshal(w.Body.Bytes(), &gotErr)

			// then
			assert.Equal(t, tt.wantCode, w.Code)

			if errBodyErr != nil {
				assert.Equal(t, tt.wantBodyErr, gotErr)
				return
			}
		})
	}
}
