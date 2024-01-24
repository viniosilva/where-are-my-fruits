package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/viniosilva/where-are-my-fruits/api/middlewares"
	"github.com/viniosilva/where-are-my-fruits/docs"
	"go.uber.org/zap"
)

//go:generate mockgen -source=./server.go -destination=../mocks/server_mocks.go -package=mocks
type HealthController interface {
	Check(ctx *gin.Context)
}

// @title			Where are my fruits API
// @version			0.0.1
// @description		Gerenciamento de frutas em baldes
// @contact.name	API Support
// @contact.email	support@wherearemyfruits.com.br
func ConfigGin(host, port string, logger *zap.SugaredLogger, health HealthController) *gin.Engine {
	r := gin.New()
	r.Use(middlewares.JSONLogMiddleware(logger))
	r.Use(middlewares.CORSMiddleware())
	r.Use(gin.Recovery())

	docs.SwaggerInfo.BasePath = "/api"
	if host != "" {
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", host, port)
	}

	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/api/healthcheck", health.Check)

	return r
}

func ConfigServer(host, port string, logger *zap.SugaredLogger, health HealthController) *http.Server {
	r := ConfigGin(host, port, logger, health)
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: r,
	}

	return server
}

// Refers: https://gin-gonic.com/docs/examples/using-middleware
//		   https://gin-gonic.com/docs/examples/grouping-routes
//		   https://github.com/swaggo/swag?tab=readme-ov-file#how-to-use-it-with-gin
