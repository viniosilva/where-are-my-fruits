package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Fruit struct {
	ID        int64      `gorm:"column:id"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`

	Name      string          `gorm:"column:name"`
	Price     decimal.Decimal `gorm:"column:price"`
	ExpiresAt time.Time       `gorm:"column:expires_at"`

	BucketID *int64 `gorm:"column:bucket_fk"`
	Bucket   Bucket `gorm:"foreignKey:bucket_fk"`
}

func (Fruit) TableName() string {
	return "fruits"
}

// Refers: https://gorm.io/docs/conventions.html#Pluralized-Table-Name
//		   https://gorm.io/docs/conventions.html#Column-Name
//		   https://gorm.io/docs/belongs_to.html
