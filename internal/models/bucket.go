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
