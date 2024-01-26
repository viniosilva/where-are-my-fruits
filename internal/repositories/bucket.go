package repositories

import "github.com/viniosilva/where-are-my-fruits/internal/models"

type BucketRepository struct {
	db DB
}

func NewBucket(db DB) *BucketRepository {
	return &BucketRepository{
		db: db,
	}
}

func (impl *BucketRepository) Create(data *models.Bucket) error {
	res := impl.db.Create(data)
	return res.Error
}
