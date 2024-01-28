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

type BucketService struct {
	db       *infra.Database
	logger   Logger
	validate Validate
}

func NewBucket(db *infra.Database, logger Logger, validate Validate) *BucketService {
	return &BucketService{
		db:       db,
		logger:   logger,
		validate: validate,
	}
}

func (impl *BucketService) Create(ctx context.Context, data dtos.CreateBucketDto) (*models.Bucket, error) {
	if err := impl.validate.Struct(data); err != nil {
		return nil, exceptions.NewValidationException(err)
	}

	bucket := models.Bucket{
		CreatedAt: _time.Now(),
		Name:      data.Name,
		Capacity:  data.Capacity,
	}

	res := impl.db.DB.Create(&bucket)
	if err := res.Error; err != nil {
		impl.logger.Error(err.Error())
		return nil, err
	}

	return &bucket, nil
}

func (impl *BucketService) List(ctx context.Context, page, pageSize int) ([]models.BucketFruits, error) {
	offset := (page - 1) * pageSize

	rows, err := impl.db.DB.Model(&models.Bucket{}).
		Select(`buckets.id,
				buckets.name,
				buckets.capacity,
				COUNT(fruits.id) AS total_fruits,
				IFNULL(SUM(fruits.price), 0) AS total_price,
				(COUNT(fruits.id) * 100 / buckets.capacity) AS percent`).
		Joins(`LEFT JOIN fruits ON fruits.bucket_fk = buckets.id
				AND fruits.deleted_at IS NULL
				AND fruits.expires_at > ?`, _time.Now()).
		Group("buckets.id").
		Order("percent DESC, buckets.created_at").
		Offset(offset).
		Limit(pageSize).
		Rows()

	if err != nil {
		impl.logger.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	bucketsFruits := make([]models.BucketFruits, 0)

	for rows.Next() {
		bucketFruits := models.BucketFruits{}
		dest := []interface{}{
			&bucketFruits.ID,
			&bucketFruits.Name,
			&bucketFruits.Capacity,
			&bucketFruits.TotalFruits,
			&bucketFruits.TotalPrice,
			&bucketFruits.Percent,
		}

		if err := rows.Scan(dest...); err != nil {
			impl.logger.Error(err.Error())
			return nil, err
		}

		bucketsFruits = append(bucketsFruits, bucketFruits)
	}

	return bucketsFruits, nil
}

func (impl *BucketService) Delete(ctx context.Context, id int64) error {
	now := _time.Now()
	err := impl.db.DB.Transaction(func(tx *gorm.DB) error {
		// Get total valid fruits by bucket
		var totalFruits int64
		res := tx.Model(&models.Fruit{}).
			Where(`bucket_fk = ?
				AND deleted_at IS NULL
				AND expires_at > ?
			`, id, now).
			Count(&totalFruits)
		if err := res.Error; err != nil {
			return err
		}
		if totalFruits > 0 {
			return exceptions.NewForbiddenException("Bucket is not empty")
		}

		return tx.Model(&models.Bucket{}).
			Where("id = ? AND deleted_at IS NULL", id).
			Update("deleted_at", now).Error
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted})

	if err != nil {
		if _, ok := err.(*exceptions.ForbiddenException); ok {
			impl.logger.Warn(err.Error())
		} else {
			impl.logger.Error(err.Error())
		}

		return err
	}

	return nil
}

// Refers: https://gorm.io/docs/scopes.html#Pagination
