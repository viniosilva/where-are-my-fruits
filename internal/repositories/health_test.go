package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/where-are-my-fruits/mocks"
)

func TestHealthRepository_NewHealth(t *testing.T) {
	t.Run("should be success", func(t *testing.T) {
		//setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sqlMock := mocks.NewMockSQL(ctrl)

		// given
		got := NewHealth(sqlMock)

		// then
		assert.NotNil(t, got)
	})
}

func TestHealthRepository_Ping(t *testing.T) {
	tests := map[string]struct {
		mock    func(sql *mocks.MockSQL)
		wantErr string
	}{
		"should be successful": {
			mock: func(sql *mocks.MockSQL) {
				sql.EXPECT().PingContext(gomock.Any()).Return(nil)
			},
		},
		"should throw error": {
			mock: func(sql *mocks.MockSQL) {
				sql.EXPECT().PingContext(gomock.Any()).Return(fmt.Errorf("error"))
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

			sqlMock := mocks.NewMockSQL(ctrl)
			tt.mock(sqlMock)

			// given
			repository := NewHealth(sqlMock)

			// when
			err := repository.Ping(ctx)

			// then
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
