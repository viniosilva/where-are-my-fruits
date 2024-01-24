package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/where-are-my-fruits/mocks"
)

func TestHealthService_NewHealth(t *testing.T) {
	t.Run("should be success", func(t *testing.T) {
		//setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		healthRepositoryMock := mocks.NewMockHealthRepository(ctrl)
		loggerMock := mocks.NewMockHealthLogger(ctrl)

		// given
		got := NewHealth(healthRepositoryMock, loggerMock)

		assert.NotNil(t, got)
	})
}

func TestHealthService_Check(t *testing.T) {
	tests := map[string]struct {
		mock    func(healthRepository *mocks.MockHealthRepository, logger *mocks.MockHealthLogger)
		wantErr string
	}{
		"should be success": {
			mock: func(healthRepository *mocks.MockHealthRepository, logger *mocks.MockHealthLogger) {
				healthRepository.EXPECT().Ping(gomock.Any()).Return(nil)
			},
		},
		"should throw error": {
			mock: func(healthRepository *mocks.MockHealthRepository, logger *mocks.MockHealthLogger) {
				healthRepository.EXPECT().Ping(gomock.Any()).Return(fmt.Errorf("error"))
				logger.EXPECT().Error(gomock.Any())
			},
			wantErr: "error",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			//setup
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			healthRepositoryMock := mocks.NewMockHealthRepository(ctrl)
			loggerMock := mocks.NewMockHealthLogger(ctrl)
			tt.mock(healthRepositoryMock, loggerMock)

			// given
			healthService := NewHealth(healthRepositoryMock, loggerMock)

			// when
			err := healthService.Check(ctx)

			// then
			if err != nil {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
