package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniosilva/where-are-my-fruits/internal/controllers/presenters"
)

type HealthController struct {
	healthService HealthService
}

//go:generate mockgen -source=./health.go -destination=../../mocks/health_controller_mocks.go -package=mocks
type HealthService interface {
	Check(ctx context.Context) error
}

func NewHealth(healthService HealthService) *HealthController {
	return &HealthController{
		healthService: healthService,
	}
}

// HealthCheck godoc
// @Summary healthcheck
// @Schemes
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} presenters.HealthCheckResponse
// @Failure 500 {object} presenters.HealthCheckResponse
// @Router /healthcheck [get]
func (impl *HealthController) Check(ctx *gin.Context) {
	code := http.StatusOK
	res := &presenters.HealthCheckResponse{
		Status: presenters.HealthCheckStatusUp,
	}

	if err := impl.healthService.Check(ctx); err != nil {
		code = http.StatusInternalServerError
		res.Status = presenters.HealthCheckStatusDown
	}

	ctx.JSON(code, res)
}
