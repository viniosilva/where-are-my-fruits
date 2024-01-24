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
		mock     func(healthService *mocks.MockHealthService)
		wantCode int
		wantBody presenters.HealthCheckResponse
	}{
		"should be success": {
			mock: func(healthService *mocks.MockHealthService) {
				healthService.EXPECT().Check(gomock.Any()).Return(nil)
			},
			wantCode: http.StatusOK,
			wantBody: presenters.HealthCheckResponse{Status: presenters.HealthCheckStatusUp},
		},
		"should throw error": {
			mock: func(healthService *mocks.MockHealthService) {
				healthService.EXPECT().Check(gomock.Any()).Return(fmt.Errorf("error"))
			},
			wantCode: http.StatusInternalServerError,
			wantBody: presenters.HealthCheckResponse{Status: presenters.HealthCheckStatusDown},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			healthServiceMock := mocks.NewMockHealthService(ctrl)
			tt.mock(healthServiceMock)

			r := gin.Default()
			healthController := NewHealth(healthServiceMock)

			r.GET("/api/healthcheck", healthController.Check)

			var got presenters.HealthCheckResponse

			// given
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/healthcheck", nil)

			// when
			r.ServeHTTP(w, req)

			json.Unmarshal(w.Body.Bytes(), &got)

			// then
			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, tt.wantBody, got)
		})
	}
}
