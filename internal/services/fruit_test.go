package services

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/viniosilva/where-are-my-fruits/internal/dtos"
	"github.com/viniosilva/where-are-my-fruits/internal/infra"
	"github.com/viniosilva/where-are-my-fruits/internal/models"
	"github.com/viniosilva/where-are-my-fruits/mocks"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestFruitService_NewFruit(t *testing.T) {
	t.Run("should be success", func(t *testing.T) {
		//setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		loggerMock := mocks.NewMockLogger(ctrl)
		validate := infra.NewValidator()

		// given
		got := NewFruit(nil, loggerMock, validate)

		// then
		assert.NotNil(t, got)
	})
}

func TestFruitService_Create(t *testing.T) {
	now := time.Date(2000, 12, 31, 23, 59, 59, 0, time.Local)
	expiresIn, _ := time.ParseDuration("1s")
	bucketID := int64(1)

	tests := map[string]struct {
		mock    func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime)
		data    dtos.CreateFruitDto
		want    *models.Fruit
		wantErr string
	}{
		"should be success when name is 128 length, price is 0 and expires in 1 second": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)
				db.ExpectBegin()
				db.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
				db.ExpectCommit()
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
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().AnyTimes().Return(now)

				bucketRows := sqlmock.NewRows([]string{"id", "name", "capacity"}).
					AddRow(int64(1), "Testing", 1)

				countTotalFruitsRows := sqlmock.NewRows([]string{"total"}).AddRow(int64(0))

				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnRows(bucketRows)               // find bucket
				db.ExpectQuery("SELECT").WillReturnRows(countTotalFruitsRows)     // count fruits per bucket
				db.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1)) // create fruit with bucketID
				db.ExpectCommit()
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
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {},
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
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {},
			data: dtos.CreateFruitDto{
				Name:      "",
				Price:     decimal.NewFromFloat32(1.99),
				ExpiresIn: &expiresIn,
			},
			wantErr: "Key: 'CreateFruitDto.Name' Error:Field validation for 'Name' failed on the 'required' tag",
		},
		"should throw error when bucket not found": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().AnyTimes().Return(now)
				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnError(fmt.Errorf(infra.MYSQL_ERROR_NOT_FOUND)) // find bucket
				db.ExpectRollback()
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
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().AnyTimes().Return(now)

				bucketRows := sqlmock.NewRows([]string{"id", "name", "capacity"}).
					AddRow(int64(1), "Testing", 1)

				countTotalFruitsRows := sqlmock.NewRows([]string{"total"}).AddRow(int64(1))

				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnRows(bucketRows)           // find bucket
				db.ExpectQuery("SELECT").WillReturnRows(countTotalFruitsRows) // count fruits per bucket
				db.ExpectRollback()

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
		"should throw error when insert": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().AnyTimes().Return(now)
				db.ExpectBegin()
				db.ExpectExec("INSERT").WillReturnError(fmt.Errorf("error")) // create fruit
				db.ExpectRollback()

				logger.EXPECT().Error(gomock.Any())
			},
			data: dtos.CreateFruitDto{
				Name:      "Testing",
				Price:     decimal.NewFromInt32(0),
				ExpiresIn: &expiresIn,
			},
			wantErr: "error",
		},
		"should throw error when bucketID is setted on insert": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().AnyTimes().Return(now)

				bucketRows := sqlmock.NewRows([]string{"id", "name", "capacity"}).
					AddRow(int64(1), "Testing", 1)

				countTotalFruitsRows := sqlmock.NewRows([]string{"total"}).AddRow(int64(0))

				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnRows(bucketRows)           // find bucket
				db.ExpectQuery("SELECT").WillReturnRows(countTotalFruitsRows) // count fruits per bucket
				db.ExpectExec("INSERT").WillReturnError(fmt.Errorf("error"))  // create fruit with bucketID
				db.ExpectRollback()

				logger.EXPECT().Error(gomock.Any())
			},
			data: dtos.CreateFruitDto{
				Name:      "Testing",
				Price:     decimal.NewFromInt32(0),
				ExpiresIn: &expiresIn,
				BucketID:  &bucketID,
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

			db, sqlMock, err := sqlmock.New()
			require.Nil(t, err)
			defer db.Close()

			dialector := mysql.New(mysql.Config{
				DSN:                       "sqlmock_db_0",
				DriverName:                "mysql",
				Conn:                      db,
				SkipInitializeWithVersion: true,
			})

			gormDB, err := gorm.Open(dialector, &gorm.Config{})
			require.Nil(t, err)
			database := &infra.Database{DB: gormDB, SQL: db}

			loggerMock := mocks.NewMockLogger(ctrl)
			validate := infra.NewValidator()
			timeMock := mocks.NewMockTime(ctrl)
			_time = timeMock

			tt.mock(sqlMock, loggerMock, timeMock)

			// given
			service := NewFruit(database, loggerMock, validate)

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
	now := time.Date(2000, 12, 31, 23, 59, 59, 0, time.Local)

	tests := map[string]struct {
		mock     func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime)
		fruitID  int64
		bucketID int64
		wantErr  string
	}{
		"should be success when bucket exists and is not full": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().AnyTimes().Return(now)

				bucketRows := sqlmock.NewRows([]string{"id", "name", "capacity"}).
					AddRow(int64(1), "Testing", 1)

				countTotalFruitsRows := sqlmock.NewRows([]string{"total"}).AddRow(int64(0))

				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnRows(bucketRows)               // find bucket
				db.ExpectQuery("SELECT").WillReturnRows(countTotalFruitsRows)     // count fruits per bucket
				db.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1)) // update fruit with bucketID
				db.ExpectCommit()
			},
			fruitID:  1,
			bucketID: 1,
		},
		"should throw error when bucket not found": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)
				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnError(fmt.Errorf(infra.MYSQL_ERROR_NOT_FOUND)) // find bucket
				db.ExpectRollback()
				logger.EXPECT().Warn(gomock.Any())
			},
			fruitID:  1,
			bucketID: 1,
			wantErr:  "Bucket not found",
		},
		"should throw error when bucket is full": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)

				bucketRows := sqlmock.NewRows([]string{"id", "name", "capacity"}).
					AddRow(int64(1), "Testing", 1)

				countTotalFruitsRows := sqlmock.NewRows([]string{"total"}).AddRow(int64(1))

				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnRows(bucketRows)           // find bucket
				db.ExpectQuery("SELECT").WillReturnRows(countTotalFruitsRows) // count fruits per bucket
				db.ExpectRollback()

				logger.EXPECT().Warn(gomock.Any())
			},
			fruitID:  1,
			bucketID: 1,
			wantErr:  "Bucket is full",
		},
		"should throw error when fruit not found": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)

				bucketRows := sqlmock.NewRows([]string{"id", "name", "capacity"}).
					AddRow(int64(1), "Testing", 1)

				countTotalFruitsRows := sqlmock.NewRows([]string{"total"}).AddRow(int64(0))

				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnRows(bucketRows)               // find bucket
				db.ExpectQuery("SELECT").WillReturnRows(countTotalFruitsRows)     // count fruits per bucket
				db.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 0)) // not found on update fruit
				db.ExpectRollback()

				logger.EXPECT().Warn(gomock.Any())
			},
			fruitID:  1,
			bucketID: 1,
			wantErr:  "Fruit not found",
		},
		"should throw error on count fruits": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)

				bucketRows := sqlmock.NewRows([]string{"id", "name", "capacity"}).
					AddRow(int64(1), "Testing", 1)

				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnRows(bucketRows)           // find bucket
				db.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("error")) // update fruit with bucketID
				db.ExpectRollback()

				logger.EXPECT().Error(gomock.Any())
			},
			fruitID:  1,
			bucketID: 1,
			wantErr:  "error",
		},
		"should throw error on update": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)

				bucketRows := sqlmock.NewRows([]string{"id", "name", "capacity"}).
					AddRow(int64(1), "Testing", 1)

				countTotalFruitsRows := sqlmock.NewRows([]string{"total"}).AddRow(int64(0))

				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnRows(bucketRows)           // find bucket
				db.ExpectQuery("SELECT").WillReturnRows(countTotalFruitsRows) // count fruits per bucket
				db.ExpectExec("UPDATE").WillReturnError(fmt.Errorf("error"))  // update fruit with bucketID
				db.ExpectRollback()

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

			db, sqlMock, err := sqlmock.New()
			require.Nil(t, err)
			defer db.Close()

			dialector := mysql.New(mysql.Config{
				DSN:                       "sqlmock_db_0",
				DriverName:                "mysql",
				Conn:                      db,
				SkipInitializeWithVersion: true,
			})

			gormDB, err := gorm.Open(dialector, &gorm.Config{})
			require.Nil(t, err)
			database := &infra.Database{DB: gormDB, SQL: db}

			loggerMock := mocks.NewMockLogger(ctrl)
			timeMock := mocks.NewMockTime(ctrl)
			_time = timeMock

			tt.mock(sqlMock, loggerMock, timeMock)

			// given
			service := NewFruit(database, loggerMock, nil)

			// when
			err = service.AddOnBucket(ctx, tt.fruitID, tt.bucketID)

			// then
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestFruitService_RemoveFromBucket(t *testing.T) {
	tests := map[string]struct {
		mock    func(db sqlmock.Sqlmock, logger *mocks.MockLogger)
		fruitID int64
		wantErr string
	}{
		"should be success when fruits exist": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger) {
				db.ExpectBegin()
				db.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1)) // update fruit
				db.ExpectCommit()
			},
			fruitID: 1,
		},
		"should throw error when fruit not found": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger) {
				db.ExpectBegin()
				db.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 0)) // update fruit
				db.ExpectCommit()

				logger.EXPECT().Warn(gomock.Any())
			},
			fruitID: 1,
			wantErr: "Fruit not found",
		},
		"should throw error": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger) {
				db.ExpectBegin()
				db.ExpectExec("UPDATE").WillReturnError(fmt.Errorf("error")) // update fruit
				db.ExpectRollback()

				logger.EXPECT().Error(gomock.Any())
			},
			fruitID: 1,
			wantErr: "error",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			//setup
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db, sqlMock, err := sqlmock.New()
			require.Nil(t, err)
			defer db.Close()

			dialector := mysql.New(mysql.Config{
				DSN:                       "sqlmock_db_0",
				DriverName:                "mysql",
				Conn:                      db,
				SkipInitializeWithVersion: true,
			})

			gormDB, err := gorm.Open(dialector, &gorm.Config{})
			require.Nil(t, err)
			database := &infra.Database{DB: gormDB, SQL: db}

			loggerMock := mocks.NewMockLogger(ctrl)

			tt.mock(sqlMock, loggerMock)

			// given
			service := NewFruit(database, loggerMock, nil)

			// when
			err = service.RemoveFromBucket(ctx, tt.fruitID)

			// then
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
func TestFruitService_Delete(t *testing.T) {
	now := time.Date(2000, 12, 31, 23, 59, 59, 0, time.Local)

	tests := map[string]struct {
		mock    func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime)
		fruitID int64
		wantErr string
	}{
		"should be success when fruit exists": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)
				db.ExpectBegin()
				db.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1)) // update fruit
				db.ExpectCommit()
			},
			fruitID: 1,
		},
		"should be success when fruit not exists": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)
				db.ExpectBegin()
				db.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 0)) // update fruit
				db.ExpectCommit()
			},
			fruitID: 1,
		},
		"should throw error": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)
				db.ExpectBegin()
				db.ExpectExec("UPDATE").WillReturnError(fmt.Errorf("error")) // update fruit
				db.ExpectRollback()

				logger.EXPECT().Error(gomock.Any())
			},
			fruitID: 1,
			wantErr: "error",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			//setup
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db, sqlMock, err := sqlmock.New()
			require.Nil(t, err)
			defer db.Close()

			dialector := mysql.New(mysql.Config{
				DSN:                       "sqlmock_db_0",
				DriverName:                "mysql",
				Conn:                      db,
				SkipInitializeWithVersion: true,
			})

			gormDB, err := gorm.Open(dialector, &gorm.Config{})
			require.Nil(t, err)
			database := &infra.Database{DB: gormDB, SQL: db}

			loggerMock := mocks.NewMockLogger(ctrl)
			timeMock := mocks.NewMockTime(ctrl)
			_time = timeMock

			tt.mock(sqlMock, loggerMock, timeMock)

			// given
			service := NewFruit(database, loggerMock, nil)

			// when
			err = service.Delete(ctx, tt.fruitID)

			// then
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
