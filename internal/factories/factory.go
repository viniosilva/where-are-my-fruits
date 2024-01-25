package factories

import (
	"database/sql"

	"github.com/viniosilva/where-are-my-fruits/internal/controllers"
	"github.com/viniosilva/where-are-my-fruits/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Factory struct {
	HealthController *controllers.HealthController
	BucketController *controllers.BucketController
}

//go:generate mockgen -source=./factory.go -destination=../../mocks/factory_mocks.go -package=mocks
type FactoryDB interface {
	DB() (*sql.DB, error)
	Create(value interface{}) (tx *gorm.DB)
}

func Build(db FactoryDB, logger *zap.SugaredLogger) (Factory, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return Factory{}, err
	}

	healthService := services.NewHealth(sqlDB, logger)
	healthController := controllers.NewHealth(healthService)

	bucketService := services.NewBucket(db, logger)
	bucketController := controllers.NewBucket(bucketService)

	return Factory{
		HealthController: healthController,
		BucketController: bucketController,
	}, nil
}
