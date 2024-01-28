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
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/where-are-my-fruits/internal/controllers/presenters"
	"github.com/viniosilva/where-are-my-fruits/internal/dtos"
	"github.com/viniosilva/where-are-my-fruits/internal/exceptions"
	"github.com/viniosilva/where-are-my-fruits/internal/models"
	"github.com/viniosilva/where-are-my-fruits/mocks"
)

func TestBucketController_Create(t *testing.T) {
	now := time.Date(2000, 12, 31, 23, 59, 59, 0, time.Local)
	tests := map[string]struct {
		mock        func(service *mocks.MockBucketService)
		body        presenters.CreateBucketReq
		wantCode    int
		wantBody    presenters.BucketRes
		wantBodyErr presenters.ErrorRes
	}{
		"should be success": {
			mock: func(service *mocks.MockBucketService) {
				data := dtos.CreateBucketDto{
					Name:     "Testing",
					Capacity: 1,
				}
				service.EXPECT().Create(gomock.Any(), data).Return(&models.Bucket{
					ID:        1,
					CreatedAt: now,
					Name:      "Testing",
					Capacity:  1,
				}, nil)
			},
			body: presenters.CreateBucketReq{
				Name:     "Testing",
				Capacity: 1,
			},
			wantCode: http.StatusCreated,
			wantBody: presenters.BucketRes{
				ID:        1,
				CreatedAt: "2000-12-31 23:59:59",
				Name:      "Testing",
				Capacity:  1,
			},
		},
		"should throw validation exception": {
			mock: func(service *mocks.MockBucketService) {
				service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, exceptions.NewValidationException(validator.ValidationErrors{
					&mocks.FieldError{Itag: "error 1", Ins: "error 1"},
					&mocks.FieldError{Itag: "error 2", Ins: "error 2"},
				}))
			},
			body:     presenters.CreateBucketReq{},
			wantCode: http.StatusBadRequest,
			wantBodyErr: presenters.ErrorRes{
				Error: exceptions.ValidationExceptionName,
				Messages: []string{
					"Key: 'error 1' Error:Field validation for '' failed on the 'error 1' tag",
					"Key: 'error 2' Error:Field validation for '' failed on the 'error 2' tag",
				},
			},
		},
		"should throw internal server error": {
			mock: func(service *mocks.MockBucketService) {
				service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			body:        presenters.CreateBucketReq{},
			wantCode:    http.StatusInternalServerError,
			wantBodyErr: presenters.ErrorRes{Error: http.StatusText(http.StatusInternalServerError)},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			serviceMock := mocks.NewMockBucketService(ctrl)
			tt.mock(serviceMock)

			r := gin.Default()
			controller := NewBucket(serviceMock)

			path := "/api/v1/buckets"
			r.POST(path, controller.Create)

			var got presenters.BucketRes
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

func TestBucketController_Delete(t *testing.T) {
	tests := map[string]struct {
		mock          func(service *mocks.MockBucketService)
		bucketIDParam string
		wantCode      int
		wantBodyErr   presenters.ErrorRes
	}{
		"should be success": {
			mock: func(service *mocks.MockBucketService) {
				service.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil)
			},
			bucketIDParam: "1",
			wantCode:      http.StatusOK,
		},
		"should throw validation exception when bucketID is invalid": {
			mock:          func(service *mocks.MockBucketService) {},
			bucketIDParam: "invalid",
			wantCode:      http.StatusBadRequest,
			wantBodyErr: presenters.ErrorRes{
				Error:   exceptions.ValidationExceptionName,
				Message: "invalid fruitID",
			},
		},
		"should throw forbidden exception when bucket is not empty": {
			mock: func(service *mocks.MockBucketService) {
				service.EXPECT().Delete(gomock.Any(), int64(1)).Return(exceptions.NewForbiddenException("Bucket is not empty"))
			},
			bucketIDParam: "1",
			wantCode:      http.StatusBadRequest,
			wantBodyErr: presenters.ErrorRes{
				Error:   exceptions.ForbiddenExceptionName,
				Message: "Bucket is not empty",
			},
		},
		"should throw internal server error": {
			mock: func(service *mocks.MockBucketService) {
				service.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
			},
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
			serviceMock := mocks.NewMockBucketService(ctrl)
			tt.mock(serviceMock)

			r := gin.Default()
			controller := NewBucket(serviceMock)

			r.DELETE("/api/v1/buckets/:bucketID", controller.Delete)

			var gotErr presenters.ErrorRes

			// given
			path := fmt.Sprintf("/api/v1/buckets/%s", tt.bucketIDParam)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", path, nil)

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
