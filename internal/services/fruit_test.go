package services

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/where-are-my-fruits/internal/dtos"
	"github.com/viniosilva/where-are-my-fruits/internal/exceptions"
	"github.com/viniosilva/where-are-my-fruits/internal/infra"
	"github.com/viniosilva/where-are-my-fruits/internal/models"
	"github.com/viniosilva/where-are-my-fruits/mocks"
)

func TestFruitService_NewFruit(t *testing.T) {
	t.Run("should be success", func(t *testing.T) {
		//setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryMock := mocks.NewMockFruitRepository(ctrl)
		loggerMock := mocks.NewMockLogger(ctrl)
		validate := infra.NewValidator()

		// given
		got := NewFruit(repositoryMock, loggerMock, validate)

		// then
		assert.NotNil(t, got)
	})
}

func TestFruitService_Create(t *testing.T) {
	now := time.Date(2000, 12, 31, 23, 59, 59, 0, time.Local)
	expiresIn, _ := time.ParseDuration("1s")
	bucketID := int64(1)

	tests := map[string]struct {
		mock    func(repository *mocks.MockFruitRepository, logger *mocks.MockLogger, time *mocks.MockTime)
		data    dtos.CreateFruitDto
		want    *models.Fruit
		wantErr string
	}{
		"should be success when name is 128 length, price is 0 and expires in 1 second": {
			mock: func(repository *mocks.MockFruitRepository, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)
				repository.EXPECT().Create(gomock.Any()).DoAndReturn(func(arg0 *models.Fruit) error {
					arg0.ID = 1
					return nil
				})
			},
			data: dtos.CreateFruitDto{
				Name:      "Testing lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris at ligula metus. Nullam eget viverra enim. Integer a vel",
				Price:     decimal.NewFromInt32(0),
				ExpiresIn: &expiresIn,
			},
			want: &models.Fruit{
				ID:        1,
				CreatedAt: now,
				Name:      "Testing lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris at ligula metus. Nullam eget viverra enim. Integer a vel",
				Price:     decimal.NewFromInt32(0),
				ExpiresAt: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local),
			},
		},
		"should be success when bucketID is setted": {
			mock: func(repository *mocks.MockFruitRepository, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)
				repository.EXPECT().Create(gomock.Any()).DoAndReturn(func(arg0 *models.Fruit) error {
					arg0.ID = 1
					return nil
				})
			},
			data: dtos.CreateFruitDto{
				Name:      "Testing",
				Price:     decimal.NewFromInt32(1),
				ExpiresIn: &expiresIn,
				BucketID:  &bucketID,
			},
			want: &models.Fruit{
				ID:        1,
				CreatedAt: now,
				Name:      "Testing",
				Price:     decimal.NewFromInt32(1),
				ExpiresAt: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local),
				BucketID:  &bucketID,
			},
		},
		"should throw error on validate when name is greater than 128, price is lower than 0 and expiresIn is empty": {
			mock: func(repository *mocks.MockFruitRepository, logger *mocks.MockLogger, mTime *mocks.MockTime) {},
			data: dtos.CreateFruitDto{
				Name:      "Testing lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris at ligula metus. Nullam eget viverra enim. Integer a veli",
				Price:     decimal.NewFromInt32(-1),
				ExpiresIn: nil,
			},
			wantErr: strings.Join([]string{
				"Key: 'CreateFruitDto.Name' Error:Field validation for 'Name' failed on the 'lte' tag",
				"Key: 'CreateFruitDto.Price' Error:Field validation for 'Price' failed on the 'dgte' tag",
				"Key: 'CreateFruitDto.ExpiresIn' Error:Field validation for 'ExpiresIn' failed on the 'required' tag",
			}, ", "),
		},
		"should throw error on validate when name is empty": {
			mock: func(repository *mocks.MockFruitRepository, logger *mocks.MockLogger, mTime *mocks.MockTime) {},
			data: dtos.CreateFruitDto{
				Name:      "",
				Price:     decimal.NewFromFloat32(1.99),
				ExpiresIn: &expiresIn,
			},
			wantErr: "Key: 'CreateFruitDto.Name' Error:Field validation for 'Name' failed on the 'required' tag",
		},
		"should throw error when bucket not found": {
			mock: func(repository *mocks.MockFruitRepository, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)
				repository.EXPECT().Create(gomock.Any()).Return(exceptions.NewForeignNotFoundException("Bucket not found"))
				logger.EXPECT().Warn(gomock.Any())
			},
			data: dtos.CreateFruitDto{
				Name:      "Testing",
				Price:     decimal.NewFromFloat32(1.99),
				ExpiresIn: &expiresIn,
				BucketID:  &bucketID,
			},
			wantErr: "Bucket not found",
		},
		"should throw error when bucket is full": {
			mock: func(repository *mocks.MockFruitRepository, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)
				repository.EXPECT().Create(gomock.Any()).Return(exceptions.NewForbiddenException("Bucket is full"))
				logger.EXPECT().Warn(gomock.Any())
			},
			data: dtos.CreateFruitDto{
				Name:      "Testing",
				Price:     decimal.NewFromFloat32(1.99),
				ExpiresIn: &expiresIn,
				BucketID:  &bucketID,
			},
			wantErr: "Bucket is full",
		},
		"should throw error": {
			mock: func(repository *mocks.MockFruitRepository, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)
				repository.EXPECT().Create(gomock.Any()).Return(&mysql.MySQLError{Message: "error"})
				logger.EXPECT().Error(gomock.Any())
			},
			data: dtos.CreateFruitDto{
				Name:      "Testing",
				Price:     decimal.NewFromInt32(0),
				ExpiresIn: &expiresIn,
			},
			wantErr: "Error 0: error",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			//setup
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repositoryMock := mocks.NewMockFruitRepository(ctrl)
			loggerMock := mocks.NewMockLogger(ctrl)
			validate := infra.NewValidator()
			timeMock := mocks.NewMockTime(ctrl)
			_time = timeMock

			tt.mock(repositoryMock, loggerMock, timeMock)

			// given
			service := NewFruit(repositoryMock, loggerMock, validate)

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

func TestFruitService_AddOnBucket(t *testing.T) {
	tests := map[string]struct {
		mock     func(repository *mocks.MockFruitRepository, logger *mocks.MockLogger, time *mocks.MockTime)
		fruitID  int64
		bucketID int64
		wantErr  string
	}{
		"should be success when bucket exists and is not full": {
			mock: func(repository *mocks.MockFruitRepository, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				repository.EXPECT().AddOnBucket(gomock.Any(), gomock.Any()).Return(nil)
			},
			fruitID:  1,
			bucketID: 1,
		},
		"should throw error when bucket not found": {
			mock: func(repository *mocks.MockFruitRepository, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				repository.EXPECT().AddOnBucket(gomock.Any(), gomock.Any()).Return(exceptions.NewForeignNotFoundException("Bucket not found"))
				logger.EXPECT().Warn(gomock.Any())
			},
			fruitID:  1,
			bucketID: 1,
			wantErr:  "Bucket not found",
		},
		"should throw error when bucket is full": {
			mock: func(repository *mocks.MockFruitRepository, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				repository.EXPECT().AddOnBucket(gomock.Any(), gomock.Any()).Return(exceptions.NewForbiddenException("Bucket is full"))
				logger.EXPECT().Warn(gomock.Any())
			},
			fruitID:  1,
			bucketID: 1,
			wantErr:  "Bucket is full",
		},
		"should throw error when fruit not found": {
			mock: func(repository *mocks.MockFruitRepository, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				repository.EXPECT().AddOnBucket(gomock.Any(), gomock.Any()).Return(exceptions.NewNotFoundException("Fruit not found"))
				logger.EXPECT().Warn(gomock.Any())
			},
			fruitID:  1,
			bucketID: 1,
			wantErr:  "Fruit not found",
		},
		"should throw error": {
			mock: func(repository *mocks.MockFruitRepository, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				repository.EXPECT().AddOnBucket(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
				logger.EXPECT().Error(gomock.Any())
			},
			fruitID:  1,
			bucketID: 1,
			wantErr:  "error",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			//setup
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repositoryMock := mocks.NewMockFruitRepository(ctrl)
			loggerMock := mocks.NewMockLogger(ctrl)

			tt.mock(repositoryMock, loggerMock, nil)

			// given
			service := NewFruit(repositoryMock, loggerMock, nil)

			// when
			err := service.AddOnBucket(ctx, tt.fruitID, tt.bucketID)

			// then
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
