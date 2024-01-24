package repositories

import (
	"context"
)

//go:generate mockgen -source=./health.go -destination=../../mocks/health_repository_mocks.go -package=mocks
type HealthDB interface {
	PingContext(ctx context.Context) error
}

type HealthRepository struct {
	db HealthDB
}

func NewHealth(db HealthDB) *HealthRepository {
	return &HealthRepository{
		db: db,
	}
}

func (impl *HealthRepository) Ping(ctx context.Context) error {
	return impl.db.PingContext(ctx)
}
