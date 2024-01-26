package repositories

import (
	"context"
)

type HealthRepository struct {
	sql SQL
}

func NewHealth(sql SQL) *HealthRepository {
	return &HealthRepository{
		sql: sql,
	}
}

func (impl *HealthRepository) Ping(ctx context.Context) error {
	return impl.sql.PingContext(ctx)
}
