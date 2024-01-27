package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/viniosilva/where-are-my-fruits/internal/infra"
	"github.com/viniosilva/where-are-my-fruits/mocks"
)

func TestHealthService_NewHealth(t *testing.T) {
	t.Run("should be success", func(t *testing.T) {
		//setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		loggerMock := mocks.NewMockLogger(ctrl)

		// given
		got := NewHealth(nil, loggerMock)

		// then
		assert.NotNil(t, got)
	})
}

func TestHealthService_Check(t *testing.T) {
	tests := map[string]struct {
		mock    func(db sqlmock.Sqlmock, logger *mocks.MockLogger)
		wantErr string
	}{
		"should be success": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger) {
				db.ExpectPing()
			},
		},
		"should throw error": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger) {
				db.ExpectPing().WillReturnError(fmt.Errorf("error"))
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

			db, sqlMock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			require.Nil(t, err)
			defer db.Close()
			database := &infra.Database{DB: nil, SQL: db}

			loggerMock := mocks.NewMockLogger(ctrl)
			tt.mock(sqlMock, loggerMock)

			// given
			service := NewHealth(database, loggerMock)

			// when
			err = service.Check(ctx)

			// then
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
