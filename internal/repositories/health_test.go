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

		// given
		dbMock := mocks.NewMockHealthDB(ctrl)

		// when
		got := NewHealth(dbMock)

		// then
		assert.NotNil(t, got)
	})
}

func TestHealthRepository_Ping(t *testing.T) {
	tests := map[string]struct {
		mock    func(db *mocks.MockHealthDB)
		wantErr string
	}{
		"should be success": {
			mock: func(db *mocks.MockHealthDB) {
				db.EXPECT().PingContext(gomock.Any()).Return(nil)
			},
		},
		"should throw error": {
			mock: func(db *mocks.MockHealthDB) {
				db.EXPECT().PingContext(gomock.Any()).Return(fmt.Errorf("error"))
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

			dbMock := mocks.NewMockHealthDB(ctrl)
			tt.mock(dbMock)

			// given
			healthRepository := NewHealth(dbMock)

			// when
			err := healthRepository.Ping(ctx)

			// then
			if err != nil {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
