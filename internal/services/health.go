package services

import (
	"context"
)

type HealthService struct {
	healthRepository HealthRepository
	logger           Logger
}

func NewHealth(repository HealthRepository, logger Logger) *HealthService {
	return &HealthService{
		healthRepository: repository,
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
