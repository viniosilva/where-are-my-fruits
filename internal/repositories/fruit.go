package repositories

import (
	"database/sql"

	"github.com/viniosilva/where-are-my-fruits/internal/exceptions"
	"github.com/viniosilva/where-are-my-fruits/internal/models"
	"gorm.io/gorm"
)

type FruitRepository struct {
	db DB
}

func NewFruit(db DB) *FruitRepository {
	return &FruitRepository{
		db: db,
	}
}

func (impl *FruitRepository) Create(data *models.Fruit) error {
	if data.BucketID == nil {
		res := impl.db.Create(data)
		return res.Error
	}

	err := impl.db.Transaction(func(tx *gorm.DB) error {
		if err := impl.ValidateBucket(*data.BucketID, tx); err != nil {
			return err
		}

		res := tx.Create(data)
		if res.Error != nil {
			return res.Error
		}

		return nil
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted})

	return err
}

func (impl *FruitRepository) ValidateBucket(bucketID int64, db *gorm.DB) error {
	now := _time.Now()

	var bucket models.Bucket
	res := db.Where("id = ?", bucketID).First(&bucket)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return exceptions.NewForeignDoesntExistsException("Bucket doesn't exists")
	}

	var totalFruits int64
	res = db.Model(&models.Fruit{}).
		Where(`bucket_fk = ?
			AND deleted_at IS NULL
			AND expires_at > ?
		`, bucket.ID, now).
		Count(&totalFruits)
	if res.Error != nil {
		return res.Error
	}
	if totalFruits >= int64(bucket.Capacity) {
		return exceptions.NewForbiddenException("Bucket is full")
	}

	return nil
}
