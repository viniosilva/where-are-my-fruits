package factories

import (
	"github.com/go-playground/validator/v10"
	"github.com/viniosilva/where-are-my-fruits/internal/controllers"
	"github.com/viniosilva/where-are-my-fruits/internal/infra"
	"github.com/viniosilva/where-are-my-fruits/internal/repositories"
	"github.com/viniosilva/where-are-my-fruits/internal/services"
	"go.uber.org/zap"
)

type Factory struct {
	HealthController *controllers.HealthController
	BucketController *controllers.BucketController
	FruitController  *controllers.FruitController
}

func Build(database *infra.Database, logger *zap.SugaredLogger, validate *validator.Validate) (Factory, error) {
	healthRepository := repositories.NewHealth(database.SQL)
	bucketRepository := repositories.NewBucket(database.DB)
	fruitRepository := repositories.NewFruit(database.DB)

	healthService := services.NewHealth(healthRepository, logger)
	bucketService := services.NewBucket(bucketRepository, logger, validate)
	fruitService := services.NewFruit(fruitRepository, logger, validate)

	healthController := controllers.NewHealth(healthService)
	bucketController := controllers.NewBucket(bucketService)
	fruitController := controllers.NewFruit(fruitService)

	return Factory{
		HealthController: healthController,
		BucketController: bucketController,
		FruitController:  fruitController,
	}, nil
}
