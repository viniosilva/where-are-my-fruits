package services

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/viniosilva/where-are-my-fruits/internal/dtos"
	"github.com/viniosilva/where-are-my-fruits/internal/exceptions"
	"github.com/viniosilva/where-are-my-fruits/internal/helpers"
	"github.com/viniosilva/where-are-my-fruits/internal/models"
	"gorm.io/gorm"
)

var _time helpers.Time

func init() {
	_time = &helpers.TimeImpl{}
}

type BucketService struct {
	repository BucketRepository
	logger     BucketLogger
	validate   *validator.Validate
}

//go:generate mockgen -source=./bucket.go -destination=../../mocks/bucket_service_mocks.go -package=mocks
type BucketLogger interface {
	Error(args ...interface{})
}

type BucketRepository interface {
	Create(value interface{}) (tx *gorm.DB)
}

func NewBucket(repository BucketRepository, logger HealthLogger) *BucketService {
	return &BucketService{
		repository: repository,
		logger:     logger,
		validate:   validator.New(),
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

	res := impl.repository.Create(&bucket)

	if res.Error != nil {
		impl.logger.Error(res.Error.Error())
		return nil, res.Error
	}

	return &bucket, nil
}
