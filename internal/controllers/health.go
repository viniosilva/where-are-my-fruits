package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniosilva/where-are-my-fruits/internal/controllers/presenters"
)

type HealthController struct {
	service HealthService
}

func NewHealth(service HealthService) *HealthController {
	return &HealthController{
		service: service,
	}
}

// HealthCheck godoc
// @Summary healthcheck
// @Schemes
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} presenters.HealthCheckRes
// @Failure 500 {object} presenters.HealthCheckRes
// @Router /healthcheck [get]
func (impl *HealthController) Check(ctx *gin.Context) {
	code := http.StatusOK
	res := &presenters.HealthCheckRes{
		Status: presenters.HealthCheckStatusUp,
	}

	if err := impl.service.Check(ctx); err != nil {
		code = http.StatusInternalServerError
		res.Status = presenters.HealthCheckStatusDown
	}

	ctx.JSON(code, res)
}
