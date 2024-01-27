package services

import (
	"context"

	"github.com/viniosilva/where-are-my-fruits/internal/dtos"
	"github.com/viniosilva/where-are-my-fruits/internal/exceptions"
	"github.com/viniosilva/where-are-my-fruits/internal/infra"
	"github.com/viniosilva/where-are-my-fruits/internal/models"
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
