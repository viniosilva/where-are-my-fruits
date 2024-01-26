package models

import (
	"time"
)

type Bucket struct {
	ID        int64      `gorm:"column:id"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`

	Name     string `gorm:"column:name"`
	Capacity int    `gorm:"column:capacity"`
}

func (Bucket) TableName() string {
	return "buckets"
}

// Refers: https://gorm.io/docs/conventions.html#Pluralized-Table-Name
//		   https://gorm.io/docs/conventions.html#Column-Name
