package repositories

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/where-are-my-fruits/internal/models"
	"github.com/viniosilva/where-are-my-fruits/mocks"
	"gorm.io/gorm"
)

func TestBucketRepository_NewBucket(t *testing.T) {
	t.Run("should be success", func(t *testing.T) {
		//setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbMock := mocks.NewMockDB(ctrl)

		// given
		got := NewBucket(dbMock)

		// then
		assert.NotNil(t, got)
	})
}

func TestBucketRepository_Create(t *testing.T) {
	now := time.Now()

	tests := map[string]struct {
		mock    func(db *mocks.MockDB)
		data    *models.Bucket
		wantErr string
	}{
		"should be successful": {
			mock: func(db *mocks.MockDB) {
				db.EXPECT().Create(gomock.Any()).DoAndReturn(func(arg0 *models.Bucket) *gorm.DB {
					arg0.ID = 1
					return &gorm.DB{RowsAffected: 1}
				})
			},
			data: &models.Bucket{
				CreatedAt: now,
				Name:      "Testing",
				Capacity:  1,
			},
		},
		"should throw error": {
			mock: func(db *mocks.MockDB) {
				db.EXPECT().Create(gomock.Any()).Return(&gorm.DB{Error: fmt.Errorf("error")})
			},
			data: &models.Bucket{
				CreatedAt: now,
				Name:      "Testing",
				Capacity:  1,
			},
			wantErr: "error",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			//setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dbMock := mocks.NewMockDB(ctrl)
			tt.mock(dbMock)

			// given
			repository := NewBucket(dbMock)

			// when
			err := repository.Create(tt.data)

			// then
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
