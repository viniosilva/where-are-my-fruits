package factories

import (
	"database/sql"

	"github.com/viniosilva/where-are-my-fruits/internal/controllers"
	"github.com/viniosilva/where-are-my-fruits/internal/repositories"
	"github.com/viniosilva/where-are-my-fruits/internal/services"
	"go.uber.org/zap"
)

type Factory struct {
	HealthController *controllers.HealthController
}

//go:generate mockgen -source=./factory.go -destination=../../mocks/factory_mocks.go -package=mocks
type FactoryDB interface {
	DB() (*sql.DB, error)
}

func Build(db FactoryDB, logger *zap.SugaredLogger) (Factory, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return Factory{}, err
	}

	healthRepository := repositories.NewHealth(sqlDB)
	healthService := services.NewHealth(healthRepository, logger)
	healthController := controllers.NewHealth(healthService)

	return Factory{
		HealthController: healthController,
	}, nil
}
