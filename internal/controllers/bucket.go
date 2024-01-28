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

	ctx.JSON(http.StatusCreated, impl.parse(res))

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

func (impl *BucketController) parse(bucket *models.Bucket) presenters.BucketRes {
	return presenters.BucketRes{
		ID:        bucket.ID,
		CreatedAt: bucket.CreatedAt.Format(time.DateTime),
		Name:      bucket.Name,
		Capacity:  bucket.Capacity,
	}
}
