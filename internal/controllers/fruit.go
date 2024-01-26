package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/viniosilva/where-are-my-fruits/internal/controllers/presenters"
	"github.com/viniosilva/where-are-my-fruits/internal/dtos"
	"github.com/viniosilva/where-are-my-fruits/internal/exceptions"
	"github.com/viniosilva/where-are-my-fruits/internal/models"
)

type FruitController struct {
	service FruitService
}

func NewFruit(service FruitService) *FruitController {
	return &FruitController{
		service: service,
	}
}

// HealthCheck godoc
// @Summary create fruit
// @Schemes
// @Tags fruit
// @Accept json
// @Produce json
// @Param fruit body presenters.CreateFruitReq true "Fruit"
// @Success 201 {object} presenters.FruitRes
// @Failure 400 {object} presenters.ErrorRes
// @Failure 500 {object} presenters.ErrorRes
// @Router /v1/fruits [post]
func (impl *FruitController) Create(ctx *gin.Context) {
	var req presenters.CreateFruitReq
	ctx.BindJSON(&req)

	var expiresIn *time.Duration
	if v, err := time.ParseDuration(req.ExpiresIn); err == nil {
		expiresIn = &v
	}

	data := dtos.CreateFruitDto{
		Name:      req.Name,
		Price:     req.Price,
		ExpiresIn: expiresIn,
		BucketID:  req.BucketID,
	}

	res, err := impl.service.Create(ctx, data)
	if err != nil {
		if e, ok := err.(*exceptions.ValidationException); ok {
			ctx.JSON(http.StatusBadRequest, presenters.ErrorRes{Error: e.Name, Messages: e.Errors})
			return
		}
		if e, ok := err.(*exceptions.ForeignDoesntExistsException); ok {
			ctx.JSON(http.StatusBadRequest, presenters.ErrorRes{Error: e.Name, Message: e.Error()})
			return
		}
		if e, ok := err.(*exceptions.ForbiddenException); ok {
			ctx.JSON(http.StatusBadRequest, presenters.ErrorRes{Error: e.Name, Message: e.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, presenters.ErrorRes{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}

	ctx.JSON(http.StatusCreated, impl.parse(res))

}

func (impl *FruitController) parse(fruit *models.Fruit) presenters.FruitRes {
	res := presenters.FruitRes{
		ID:        fruit.ID,
		CreatedAt: fruit.CreatedAt.Format(time.DateTime),
		Name:      fruit.Name,
		Price:     fruit.Price,
		ExpiresAt: fruit.ExpiresAt.Format(time.DateTime),
	}

	if fruit.BucketID != nil {
		res.BucketID = fruit.BucketID
	}

	return res
}
