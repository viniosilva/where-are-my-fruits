package repositories

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/where-are-my-fruits/internal/models"
	"github.com/viniosilva/where-are-my-fruits/mocks"
	"gorm.io/gorm"
)

func TestFruitRepository_NewFruit(t *testing.T) {
	t.Run("should be success", func(t *testing.T) {
		//setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbMock := mocks.NewMockDB(ctrl)

		// given
		got := NewFruit(dbMock)

		// then
		assert.NotNil(t, got)
	})
}

func TestFruitRepository_Create(t *testing.T) {
	now := time.Now()
	bucketID := int64(1)

	tests := map[string]struct {
		mock    func(db *mocks.MockDB)
		data    *models.Fruit
		wantErr string
	}{
		"should be successful when bucketID is empty": {
			mock: func(db *mocks.MockDB) {
				db.EXPECT().Create(gomock.Any()).DoAndReturn(func(arg0 *models.Fruit) *gorm.DB {
					arg0.ID = 1
					return &gorm.DB{RowsAffected: 1}
				})
			},
			data: &models.Fruit{
				CreatedAt: now,
				Name:      "Testing",
				Price:     decimal.NewFromInt(1),
				ExpiresAt: now.Add(1 * time.Hour),
			},
		},
		"should be successful when bucketID is not empty": {
			mock: func(db *mocks.MockDB) {
				db.EXPECT().Transaction(gomock.Any(), gomock.Any()).Return(nil)
			},
			data: &models.Fruit{
				CreatedAt: now,
				Name:      "Testing",
				Price:     decimal.NewFromInt(1),
				ExpiresAt: now.Add(1 * time.Hour),
				BucketID:  &bucketID,
			},
		},
		"should throw error when create fails": {
			mock: func(db *mocks.MockDB) {
				db.EXPECT().Create(gomock.Any()).Return(&gorm.DB{Error: fmt.Errorf("error")})
			},
			data: &models.Fruit{
				CreatedAt: now,
				Name:      "Testing",
				Price:     decimal.NewFromInt(1),
				ExpiresAt: now.Add(1 * time.Hour),
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
			repository := NewFruit(dbMock)

			// when
			err := repository.Create(tt.data)

			// then
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
