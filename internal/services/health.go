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
	Ping(ctx context.Context) error
}

func NewHealth(healthRepository HealthRepository, logger HealthLogger) *HealthService {
	return &HealthService{
		healthRepository: healthRepository,
		logger:           logger,
	}
}

func (impl *HealthService) Check(ctx context.Context) error {
	err := impl.healthRepository.Ping(ctx)
	if err != nil {
		impl.logger.Error(err.Error())
	}

	return err
}
