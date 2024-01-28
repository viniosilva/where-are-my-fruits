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

type BucketController struct {
	service BucketService
}

func NewBucket(service BucketService) *BucketController {
	return &BucketController{
		service: service,
	}
}

// Bucket godoc
// @Summary create bucket
// @Schemes
// @Tags bucket
// @Accept json
// @Produce json
// @Param bucket body presenters.CreateBucketReq true "Bucket"
// @Success 201 {object} presenters.BucketRes
// @Failure 400 {object} presenters.ErrorRes
// @Failure 500 {object} presenters.ErrorRes
// @Router /v1/buckets [post]
func (impl *BucketController) Create(ctx *gin.Context) {
	var req presenters.CreateBucketReq
	ctx.BindJSON(&req)

	data := dtos.CreateBucketDto{
		Name:     req.Name,
		Capacity: req.Capacity,
	}

	res, err := impl.service.Create(ctx, data)
	if err != nil {
		if e, ok := err.(*exceptions.ValidationException); ok {
			ctx.JSON(http.StatusBadRequest, presenters.ErrorRes{Error: e.Name, Messages: e.Errors})
			return
		}

		ctx.JSON(http.StatusInternalServerError, presenters.ErrorRes{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}

	ctx.JSON(http.StatusCreated, impl.parseModel(res))
}

// Bucket godoc
// @Summary list buckets
// @Schemes
// @Tags bucket
// @Accept json
// @Produce json
// @Param page query int false "page" default(1)
// @Param pageSize query int false "pageSize" default(10)
// @Success 200 {object} presenters.BucketsRes
// @Failure 500 {object} presenters.ErrorRes
// @Router /v1/buckets [get]
func (impl *BucketController) List(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	res, err := impl.service.List(ctx, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, presenters.ErrorRes{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}

	resp := presenters.BucketsFruitsRes{Data: []presenters.BucketFruitsRes{}}
	for _, bucket := range res {
		resp.Data = append(resp.Data, impl.parseDTO(&bucket))
	}

	ctx.JSON(http.StatusOK, resp)
}

// Fruit godoc
// @Summary delete bucket
// @Schemes
// @Tags bucket
// @Accept json
// @Produce json
// @Param bucketID path int64 true "Bucket ID"
// @Success 200 {object} nil
// @Failure 400 {object} presenters.ErrorRes
// @Failure 500 {object} presenters.ErrorRes
// @Router /v1/buckets/{bucketID} [delete]
func (impl *BucketController) Delete(ctx *gin.Context) {
	bucketID, err := strconv.ParseInt(ctx.Param("bucketID"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, presenters.ErrorRes{Error: exceptions.ValidationExceptionName, Message: "invalid bucketID"})
		return
	}

	err = impl.service.Delete(ctx, bucketID)
	if err != nil {
		if e, ok := err.(*exceptions.ForbiddenException); ok {
			ctx.JSON(http.StatusBadRequest, presenters.ErrorRes{Error: e.Name, Message: e.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, presenters.ErrorRes{Error: http.StatusText(http.StatusInternalServerError)})
		return
	}

	ctx.Status(http.StatusOK)
}

func (impl *BucketController) parseModel(bucket *models.Bucket) presenters.BucketRes {
	return presenters.BucketRes{
		ID:        bucket.ID,
		CreatedAt: bucket.CreatedAt.Format(time.DateTime),
		Name:      bucket.Name,
		Capacity:  bucket.Capacity,
	}
}

func (impl *BucketController) parseDTO(bucket *models.BucketFruits) presenters.BucketFruitsRes {
	return presenters.BucketFruitsRes{
		ID:          bucket.ID,
		Name:        bucket.Name,
		Capacity:    bucket.Capacity,
		TotalFruits: bucket.TotalFruits,
		TotalPrice:  bucket.TotalPrice,
		Percent:     bucket.Percent.StringFixed(2) + "%",
	}
}
