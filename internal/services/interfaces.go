package services

import (
	"context"

	"github.com/viniosilva/where-are-my-fruits/internal/helpers"
	"github.com/viniosilva/where-are-my-fruits/internal/models"
)

var _time helpers.Time

func init() {
	_time = &helpers.TimeImpl{}
}

//go:generate mockgen -source=./interfaces.go -destination=../../mocks/services_mocks.go -package=mocks
type Logger interface {
	Warn(args ...interface{})
	Error(args ...interface{})
}

type Validate interface {
	Struct(s interface{}) error
}

type HealthRepository interface {
	Ping(ctx context.Context) error
}

type BucketRepository interface {
	Create(data *models.Bucket) error
}

type FruitRepository interface {
	Create(data *models.Fruit) error
	AddOnBucket(fruitID, bucketID int64) error
}
