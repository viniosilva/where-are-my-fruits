package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/viniosilva/where-are-my-fruits/api"
	"github.com/viniosilva/where-are-my-fruits/internal/factories"
	"github.com/viniosilva/where-are-my-fruits/internal/infra"
)

func main() {
	logger := infra.ConfigLogger()

	config, err := infra.GetConfig(".")
	if err != nil {
		log.Fatalf("infra.GetConfig: %s\n", err)
	}

	db, err := infra.ConfigDB(config.MySQL.Username, config.MySQL.Password, config.MySQL.Host,
		config.MySQL.Port, config.MySQL.Database, config.MySQL.ConnMaxLifetime,
		config.MySQL.MaxIdleConns, config.MySQL.MaxOpenConns)
	if err != nil {
		log.Fatalf("infra.ConfigDB: %s\n", err)
	}

	validate := infra.NewValidator()

	factory, err := factories.Build(db, logger, validate)
	if err != nil {
		log.Fatalf("factory.Build: %s\n", err)
	}

	server := api.ConfigServer(config.Api.Host, config.Api.Port, logger, factory.HealthController, factory.BucketController, factory.FruitController)

	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("api.Run: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutdown server...")

	db.SQL.Close()

	ctx := context.Background()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server.Shutdown: %s\n", err)
	}

	logger.Info("Bye")
}

// Refers: https://gin-gonic.com/docs/examples/graceful-restart-or-stop
