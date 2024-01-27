package services

import (
	"context"

	"github.com/viniosilva/where-are-my-fruits/internal/infra"
)

type HealthService struct {
	db     *infra.Database
	logger Logger
}

func NewHealth(db *infra.Database, logger Logger) *HealthService {
	return &HealthService{
		db:     db,
		logger: logger,
	}
}

func (impl *HealthService) Check(ctx context.Context) error {
	err := impl.db.SQL.PingContext(ctx)
	if err != nil {
		impl.logger.Error(err.Error())
	}

	return err
}
