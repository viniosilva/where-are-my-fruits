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

		repositoryMock := mocks.NewMockHealthRepository(ctrl)
		loggerMock := mocks.NewMockLogger(ctrl)

		// given
		got := NewHealth(repositoryMock, loggerMock)

		// then
		assert.NotNil(t, got)
	})
}

func TestHealthService_Check(t *testing.T) {
	tests := map[string]struct {
		mock    func(repository *mocks.MockHealthRepository, logger *mocks.MockLogger)
		wantErr string
	}{
		"should be success": {
			mock: func(repository *mocks.MockHealthRepository, logger *mocks.MockLogger) {
				repository.EXPECT().Ping(gomock.Any()).Return(nil)
			},
		},
		"should throw error": {
			mock: func(repository *mocks.MockHealthRepository, logger *mocks.MockLogger) {
				repository.EXPECT().Ping(gomock.Any()).Return(fmt.Errorf("error"))
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

			repositoryMock := mocks.NewMockHealthRepository(ctrl)
			loggerMock := mocks.NewMockLogger(ctrl)
			tt.mock(repositoryMock, loggerMock)

			// given
			service := NewHealth(repositoryMock, loggerMock)

			// when
			err := service.Check(ctx)

			// then
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
