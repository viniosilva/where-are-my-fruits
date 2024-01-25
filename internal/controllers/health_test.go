package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/where-are-my-fruits/internal/controllers/presenters"
	"github.com/viniosilva/where-are-my-fruits/mocks"
)

func TestHealthController_Check(t *testing.T) {
	tests := map[string]struct {
		mock     func(service *mocks.MockHealthService)
		wantCode int
		wantBody presenters.HealthCheckRes
	}{
		"should be success": {
			mock: func(service *mocks.MockHealthService) {
				service.EXPECT().Check(gomock.Any()).Return(nil)
			},
			wantCode: http.StatusOK,
			wantBody: presenters.HealthCheckRes{Status: presenters.HealthCheckStatusUp},
		},
		"should throw error": {
			mock: func(service *mocks.MockHealthService) {
				service.EXPECT().Check(gomock.Any()).Return(fmt.Errorf("error"))
			},
			wantCode: http.StatusInternalServerError,
			wantBody: presenters.HealthCheckRes{Status: presenters.HealthCheckStatusDown},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			serviceMock := mocks.NewMockHealthService(ctrl)
			tt.mock(serviceMock)

			r := gin.Default()
			controller := NewHealth(serviceMock)

			path := "/api/healthcheck"
			r.GET(path, controller.Check)

			var got presenters.HealthCheckRes

			// given
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", path, nil)

			// when
			r.ServeHTTP(w, req)

			json.Unmarshal(w.Body.Bytes(), &got)

			// then
			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, tt.wantBody, got)
		})
	}
}
