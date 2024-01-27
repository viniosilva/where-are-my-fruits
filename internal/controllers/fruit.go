package controllers

import (
	"net/http"
	"strconv"
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

// Fruit godoc
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
		if e, ok := err.(*exceptions.ForeignNotFoundException); ok {
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

// Fruit godoc
// @Summary add fruit on bucket
// @Schemes
// @Tags fruit
// @Accept json
// @Produce json
// @Param fruitID path int64 true "Fruit ID"
// @Param bucketID path int64 true "Bucket ID"
// @Success 200 {object} nil
// @Failure 400 {object} presenters.ErrorRes
// @Failure 500 {object} presenters.ErrorRes
// @Router /v1/fruits/{fruitID}/buckets/{bucketID} [post]
func (impl *FruitController) AddOnBucket(ctx *gin.Context) {
	fruitID, err := strconv.ParseInt(ctx.Param("fruitID"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, presenters.ErrorRes{Error: exceptions.ValidationExceptionName, Message: "invalid fruitID"})
		return
	}

	bucketID, err := strconv.ParseInt(ctx.Param("bucketID"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, presenters.ErrorRes{Error: exceptions.ValidationExceptionName, Message: "invalid bucketID"})
		return
	}

	err = impl.service.AddOnBucket(ctx, fruitID, bucketID)
	if err != nil {
		if e, ok := err.(*exceptions.ForeignNotFoundException); ok {
			ctx.JSON(http.StatusBadRequest, presenters.ErrorRes{Error: e.Name, Message: e.Error()})
			return
		}
		if e, ok := err.(*exceptions.ForbiddenException); ok {
			ctx.JSON(http.StatusBadRequest, presenters.ErrorRes{Error: e.Name, Message: e.Error()})
			return
		}
		if e, ok := err.(*exceptions.NotFoundException); ok {
			ctx.JSON(http.StatusBadRequest, presenters.ErrorRes{Error: e.Name, Message: e.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, presenters.ErrorRes{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}

	ctx.Status(http.StatusOK)
}

// Fruit godoc
// @Summary remove fruit from bucket
// @Schemes
// @Tags fruit
// @Accept json
// @Produce json
// @Param fruitID path int64 true "Fruit ID"
// @Success 200 {object} nil
// @Failure 400 {object} presenters.ErrorRes
// @Failure 500 {object} presenters.ErrorRes
// @Router /v1/fruits/{fruitID}/buckets [delete]
func (impl *FruitController) RemoveFromBucket(ctx *gin.Context) {
	fruitID, err := strconv.ParseInt(ctx.Param("fruitID"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, presenters.ErrorRes{Error: exceptions.ValidationExceptionName, Message: "invalid fruitID"})
		return
	}

	err = impl.service.RemoveFromBucket(ctx, fruitID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, presenters.ErrorRes{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}

	ctx.Status(http.StatusOK)
}

// Fruit godoc
// @Summary delete fruit
// @Schemes
// @Tags fruit
// @Accept json
// @Produce json
// @Param fruitID path int64 true "Fruit ID"
// @Success 200 {object} nil
// @Failure 400 {object} presenters.ErrorRes
// @Failure 500 {object} presenters.ErrorRes
// @Router /v1/fruits/{fruitID} [delete]
func (impl *FruitController) Delete(ctx *gin.Context) {
	fruitID, err := strconv.ParseInt(ctx.Param("fruitID"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, presenters.ErrorRes{Error: exceptions.ValidationExceptionName, Message: "invalid fruitID"})
		return
	}

	err = impl.service.Delete(ctx, fruitID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, presenters.ErrorRes{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}

	ctx.Status(http.StatusOK)
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
