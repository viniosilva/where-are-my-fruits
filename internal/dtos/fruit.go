package dtos

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateFruitDto struct {
	Name      string          `validate:"required,gt=0,lte=128"`
	Price     decimal.Decimal `validate:"required,dgte=0"`
	ExpiresIn *time.Duration  `validate:"required"`
	BucketID  *int64          `validate:"omitempty,gt=0"`
}
