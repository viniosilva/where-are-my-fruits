package services

import (
	"context"
)

type HealthService struct {
	healthRepository HealthRepository
	logger           HealthLogger
}

//go:generate mockgen -source=./health.go -destination=../../mocks/health_service_mocks.go -package=mocks
type HealthLogger interface {
	Error(args ...interface{})
}

type HealthRepository interface {
	PingContext(ctx context.Context) error
}

func NewHealth(repository HealthRepository, logger HealthLogger) *HealthService {
	return &HealthService{
		healthRepository: repository,
		logger:           logger,
	}
}

func (impl *HealthService) Check(ctx context.Context) error {
	err := impl.healthRepository.PingContext(ctx)
	if err != nil {
		impl.logger.Error(err.Error())
	}

	return err
}
