package services

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/where-are-my-fruits/internal/dtos"
	"github.com/viniosilva/where-are-my-fruits/internal/infra"
	"github.com/viniosilva/where-are-my-fruits/internal/models"
	"github.com/viniosilva/where-are-my-fruits/mocks"
	"gorm.io/gorm"
)

func TestBucketService_NewBucket(t *testing.T) {
	t.Run("should be success", func(t *testing.T) {
		//setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryMock := mocks.NewMockBucketRepository(ctrl)
		loggerMock := mocks.NewMockLogger(ctrl)
		validate := infra.NewValidator()

		// given
		got := NewBucket(repositoryMock, loggerMock, validate)

		// then
		assert.NotNil(t, got)
	})
}

func TestBucketService_Create(t *testing.T) {
	now := time.Now()

	tests := map[string]struct {
		mock    func(repository *mocks.MockBucketRepository, logger *mocks.MockLogger, time *mocks.MockTime)
		data    dtos.CreateBucketDto
		want    *models.Bucket
		wantErr string
	}{
		"should be success when name is 128 length and capacity is 1": {
			mock: func(repository *mocks.MockBucketRepository, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)
				repository.EXPECT().Create(gomock.Any()).DoAndReturn(func(arg0 *models.Bucket) *gorm.DB {
					arg0.ID = 1
					return &gorm.DB{RowsAffected: 1}
				})
			},
			data: dtos.CreateBucketDto{
				Name:     "Testing lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris at ligula metus. Nullam eget viverra enim. Integer a vel",
				Capacity: 1,
			},
			want: &models.Bucket{
				ID:        1,
				CreatedAt: now,
				Name:      "Testing lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris at ligula metus. Nullam eget viverra enim. Integer a vel",
				Capacity:  1,
			},
		},
		"should throw error on validate when name is greater than 128 and capacity is lower then 1": {
			mock: func(repository *mocks.MockBucketRepository, logger *mocks.MockLogger, mTime *mocks.MockTime) {},
			data: dtos.CreateBucketDto{
				Name:     "Testing lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris at ligula metus. Nullam eget viverra enim. Integer a veli",
				Capacity: 0,
			},
			wantErr: strings.Join([]string{
				"Key: 'CreateBucketDto.Name' Error:Field validation for 'Name' failed on the 'lte' tag",
				"Key: 'CreateBucketDto.Capacity' Error:Field validation for 'Capacity' failed on the 'required' tag",
			}, ", "),
		},
		"should throw error on validate when name is empty": {
			mock: func(repository *mocks.MockBucketRepository, logger *mocks.MockLogger, mTime *mocks.MockTime) {},
			data: dtos.CreateBucketDto{
				Name:     "",
				Capacity: 1,
			},
			wantErr: "Key: 'CreateBucketDto.Name' Error:Field validation for 'Name' failed on the 'required' tag",
		},
		"should throw error": {
			mock: func(repository *mocks.MockBucketRepository, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)
				repository.EXPECT().Create(gomock.Any()).Return(fmt.Errorf("error"))
				logger.EXPECT().Error(gomock.Any())
			},
			data: dtos.CreateBucketDto{
				Name:     "Testing",
				Capacity: 1,
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

			repositoryMock := mocks.NewMockBucketRepository(ctrl)
			loggerMock := mocks.NewMockLogger(ctrl)
			validate := infra.NewValidator()
			timeMock := mocks.NewMockTime(ctrl)
			_time = timeMock

			tt.mock(repositoryMock, loggerMock, timeMock)

			// given
			service := NewBucket(repositoryMock, loggerMock, validate)

			// when
			got, err := service.Create(ctx, tt.data)

			// then
			assert.Equal(t, tt.want, got)
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
