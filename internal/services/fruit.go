package services

import (
	"context"
	"database/sql"

	"github.com/viniosilva/where-are-my-fruits/internal/dtos"
	"github.com/viniosilva/where-are-my-fruits/internal/exceptions"
	"github.com/viniosilva/where-are-my-fruits/internal/infra"
	"github.com/viniosilva/where-are-my-fruits/internal/models"
	"gorm.io/gorm"
)

type FruitService struct {
	db       *infra.Database
	logger   Logger
	validate Validate
}

func NewFruit(db *infra.Database, logger Logger, validate Validate) *FruitService {
	return &FruitService{
		db:       db,
		logger:   logger,
		validate: validate,
	}
}

func (impl *FruitService) Create(ctx context.Context, data dtos.CreateFruitDto) (*models.Fruit, error) {
	if err := impl.validate.Struct(data); err != nil {
		return nil, exceptions.NewValidationException(err)
	}

	now := _time.Now()
	fruit := models.Fruit{
		CreatedAt: now,
		Name:      data.Name,
		Price:     data.Price,
		ExpiresAt: now.Add(*data.ExpiresIn),
	}

	if data.BucketID == nil {
		res := impl.db.DB.Create(&fruit)
		if err := res.Error; err != nil {
			impl.logger.Error(err.Error())
			return nil, err
		}

		return &fruit, nil
	}

	fruit.BucketID = data.BucketID
	err := impl.db.DB.Transaction(func(tx *gorm.DB) error {
		if err := impl.validateBucket(ctx, tx, *data.BucketID); err != nil {
			return err
		}

		return tx.Create(&fruit).Error
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted})

	if err != nil {
		if _, ok := err.(*exceptions.ForeignNotFoundException); ok {
			impl.logger.Warn(err.Error())
		} else if _, ok := err.(*exceptions.ForbiddenException); ok {
			impl.logger.Warn(err.Error())
		} else {
			impl.logger.Error(err.Error())
		}

		return nil, err
	}

	return &fruit, nil
}

func (impl *FruitService) AddOnBucket(ctx context.Context, fruitID, bucketID int64) error {
	err := impl.db.DB.Transaction(func(tx *gorm.DB) error {
		err := impl.validateBucket(ctx, tx, bucketID)
		if err != nil {
			return err
		}

		res := tx.Model(&models.Fruit{}).
			Where("id = ?", fruitID).
			Update("bucket_fk", bucketID)

		if err := res.Error; err != nil {
			return err
		}
		if res.RowsAffected == 0 {
			return exceptions.NewNotFoundException("Fruit not found")
		}
		return nil
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted})

	if err != nil {
		if _, ok := err.(*exceptions.NotFoundException); ok {
			impl.logger.Warn(err.Error())
		} else if _, ok := err.(*exceptions.ForeignNotFoundException); ok {
			impl.logger.Warn(err.Error())
		} else if _, ok := err.(*exceptions.ForbiddenException); ok {
			impl.logger.Warn(err.Error())
		} else {
			impl.logger.Error(err.Error())
		}

		return err
	}

	return nil
}

func (impl *FruitService) RemoveFromBucket(ctx context.Context, fruitID int64) error {
	res := impl.db.DB.Model(&models.Fruit{}).
		Where("id = ?", fruitID).
		Update("bucket_fk", nil)

	if err := res.Error; err != nil {
		impl.logger.Error(err.Error())
		return err
	}
	if res.RowsAffected == 0 {
		err := exceptions.NewNotFoundException("Fruit not found")
		impl.logger.Warn(err.Error())
		return err
	}

	return res.Error
}

func (impl *FruitService) Delete(ctx context.Context, fruitID int64) error {
	res := impl.db.DB.Model(&models.Fruit{}).
		Where("id = ?", fruitID).
		Update("deleted_at", _time.Now())

	if err := res.Error; err != nil {
		impl.logger.Error(err.Error())
		return err
	}

	return nil
}

func (impl *FruitService) validateBucket(ctx context.Context, tx *gorm.DB, bucketID int64) error {
	now := _time.Now()

	// Get bucket by ID
	var bucket models.Bucket
	res := tx.Where("id = ?", bucketID).First(&bucket)
	err := res.Error
	if err != nil && err.Error() == infra.MYSQL_ERROR_NOT_FOUND {
		return exceptions.NewForeignNotFoundException("Bucket not found")
	}

	// Get total valid fruits by bucket
	var totalFruits int64
	res = tx.Model(&models.Fruit{}).
		Where(`bucket_fk = ?
			AND deleted_at IS NULL
			AND expires_at > ?
		`, bucketID, now).
		Count(&totalFruits)
	if err := res.Error; err != nil {
		return err
	}

	// Validate current bucket capacity
	if totalFruits >= int64(bucket.Capacity) {
		return exceptions.NewForbiddenException("Bucket is full")
	}

	return nil
}

// Refers: https://gorm.io/docs/transactions.html#Transaction
