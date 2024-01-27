package services

import (
	"context"

	"github.com/viniosilva/where-are-my-fruits/internal/dtos"
	"github.com/viniosilva/where-are-my-fruits/internal/exceptions"
	"github.com/viniosilva/where-are-my-fruits/internal/models"
)

type FruitService struct {
	repository FruitRepository
	logger     Logger
	validate   Validate
}

func NewFruit(repository FruitRepository, logger Logger, validate Validate) *FruitService {
	return &FruitService{
		repository: repository,
		logger:     logger,
		validate:   validate,
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

	if data.BucketID != nil {
		fruit.BucketID = data.BucketID
	}

	err := impl.repository.Create(&fruit)

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

	return &fruit, err
}

func (impl *FruitService) AddOnBucket(ctx context.Context, fruitID, bucketID int64) error {
	err := impl.repository.AddOnBucket(fruitID, bucketID)

	if err != nil {
		if _, ok := err.(*exceptions.ForeignNotFoundException); ok {
			impl.logger.Warn(err.Error())
		} else if _, ok := err.(*exceptions.ForbiddenException); ok {
			impl.logger.Warn(err.Error())
		} else if _, ok := err.(*exceptions.NotFoundException); ok {
			impl.logger.Warn(err.Error())
		} else {
			impl.logger.Error(err.Error())
		}

		return err
	}

	return nil
}

func (impl *FruitService) RemoveFromBucket(ctx context.Context, fruitID, bucketID int64) error {
	err := impl.repository.RemoveFromBucket(fruitID, bucketID)

	if err != nil {
		impl.logger.Error(err.Error())
		return err
	}

	return nil
}
