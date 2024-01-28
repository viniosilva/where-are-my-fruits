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

func TestBucketService_NewBucket(t *testing.T) {
	t.Run("should be success", func(t *testing.T) {
		//setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		loggerMock := mocks.NewMockLogger(ctrl)
		validate := infra.NewValidator()

		// given
		got := NewBucket(nil, loggerMock, validate)

		// then
		assert.NotNil(t, got)
	})
}

func TestBucketService_Create(t *testing.T) {
	now := time.Now()

	tests := map[string]struct {
		mock    func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime)
		data    dtos.CreateBucketDto
		want    *models.Bucket
		wantErr string
	}{
		"should be success when name is 128 length and capacity is 1": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)
				db.ExpectBegin()
				db.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
				db.ExpectCommit()
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
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {},
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
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {},
			data: dtos.CreateBucketDto{
				Name:     "",
				Capacity: 1,
			},
			wantErr: "Key: 'CreateBucketDto.Name' Error:Field validation for 'Name' failed on the 'required' tag",
		},
		"should throw error": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().AnyTimes().Return(now)
				db.ExpectBegin()
				db.ExpectExec("INSERT").WillReturnError(fmt.Errorf("error"))
				db.ExpectRollback()

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
			service := NewBucket(database, loggerMock, validate)

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

func TestBucketService_List(t *testing.T) {
	now := time.Now()

	tests := map[string]struct {
		mock     func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime)
		page     int
		pageSize int
		want     []models.BucketFruits
		wantErr  string
	}{
		"shoulb be successful": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)

				rows := sqlmock.
					NewRows([]string{"id", "name", "capacity", "total_fruits", "total_price", "percent"}).
					AddRow(int64(1), "Testing", 4, int64(3), decimal.NewFromFloat32(16.32), decimal.NewFromInt32(75)).
					AddRow(int64(2), "Testing_2", 3, int64(1), decimal.NewFromFloat32(6.25), decimal.NewFromFloat32(33.33))

				db.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			page:     1,
			pageSize: 10,
			want: []models.BucketFruits{
				{
					ID:          1,
					Name:        "Testing",
					Capacity:    4,
					TotalFruits: 3,
					TotalPrice:  decimal.NewFromFloat32(16.32),
					Percent:     decimal.NewFromInt32(75),
				},
				{
					ID:          2,
					Name:        "Testing_2",
					Capacity:    3,
					TotalFruits: 1,
					TotalPrice:  decimal.NewFromFloat32(6.25),
					Percent:     decimal.NewFromFloat32(33.33),
				},
			},
		},
		"shoulb throw error when select": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)

				db.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("error"))
				logger.EXPECT().Error(gomock.Any())
			},
			page:     1,
			pageSize: 10,
			wantErr:  "error",
		},
		"shoulb throw error when scan": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)

				rows := sqlmock.
					NewRows([]string{"id", "name", "capacity", "total_fruits", "total_price", "percent"}).
					AddRow(int64(1), "Testing", 4, int64(3), decimal.NewFromFloat32(16.32), nil)

				db.ExpectQuery("SELECT").WillReturnRows(rows)
				logger.EXPECT().Error(gomock.Any())
			},
			page:     1,
			pageSize: 10,
			wantErr:  "sql: Scan error on column index 5, name \"percent\": could not convert value '<nil>' to byte array of type '<nil>'",
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
			service := NewBucket(database, loggerMock, validate)

			// when
			got, err := service.List(ctx, tt.page, tt.pageSize)

			// then
			assert.Equal(t, tt.want, got)
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestBucketService_Delete(t *testing.T) {
	now := time.Date(2000, 12, 31, 23, 59, 59, 0, time.Local)

	tests := map[string]struct {
		mock    func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime)
		fruitID int64
		wantErr string
	}{
		"should be success when bucket exists": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)

				countTotalFruitsRows := sqlmock.NewRows([]string{"total"}).AddRow(int64(0))

				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnRows(countTotalFruitsRows)     // count fruits per bucket
				db.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1)) // update bucket
				db.ExpectCommit()
			},
			fruitID: 1,
		},
		"should be success when bucket not exists": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)

				countTotalFruitsRows := sqlmock.NewRows([]string{"total"}).AddRow(int64(0))

				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnRows(countTotalFruitsRows)     // count fruits per bucket
				db.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 0)) // update bucket
				db.ExpectCommit()
			},
			fruitID: 1,
		},
		"should throw forbidden error when bucker is not empty": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)

				countTotalFruitsRows := sqlmock.NewRows([]string{"total"}).AddRow(int64(1))

				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnRows(countTotalFruitsRows) // count fruits per bucket
				db.ExpectRollback()

				logger.EXPECT().Warn(gomock.Any())
			},
			fruitID: 1,
			wantErr: "Bucket is not empty",
		},
		"should throw error when count total fruits": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)

				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("error")) // count fruits per bucket
				db.ExpectRollback()

				logger.EXPECT().Error(gomock.Any())
			},
			fruitID: 1,
			wantErr: "error",
		},
		"should throw error on update": {
			mock: func(db sqlmock.Sqlmock, logger *mocks.MockLogger, mTime *mocks.MockTime) {
				mTime.EXPECT().Now().Return(now)

				countTotalFruitsRows := sqlmock.NewRows([]string{"total"}).AddRow(int64(0))

				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnRows(countTotalFruitsRows) // count fruits per bucket
				db.ExpectExec("UPDATE").WillReturnError(fmt.Errorf("error"))  // update bucket
				db.ExpectCommit()

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
			service := NewBucket(database, loggerMock, nil)

			// when
			err = service.Delete(ctx, tt.fruitID)

			// then
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
