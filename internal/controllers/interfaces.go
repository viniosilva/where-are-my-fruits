package controllers

import (
	"context"

	"github.com/viniosilva/where-are-my-fruits/internal/dtos"
	"github.com/viniosilva/where-are-my-fruits/internal/models"
)

//go:generate mockgen -source=./interfaces.go -destination=../../mocks/controllers_mocks.go -package=mocks
type HealthService interface {
	Check(ctx context.Context) error
}

type BucketService interface {
	Create(ctx context.Context, data dtos.CreateBucketDto) (*models.Bucket, error)
	List(ctx context.Context, page, pageSize int) ([]models.BucketFruits, error)
	Delete(ctx context.Context, id int64) error
}

type FruitService interface {
	Create(ctx context.Context, data dtos.CreateFruitDto) (*models.Fruit, error)
	AddOnBucket(ctx context.Context, fruitID, bucketID int64) error
	RemoveFromBucket(ctx context.Context, fruitID int64) error
	Delete(ctx context.Context, id int64) error
}
