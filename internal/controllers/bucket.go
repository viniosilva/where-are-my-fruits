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

func (impl *BucketController) parse(bucket *models.Bucket) presenters.BucketRes {
	return presenters.BucketRes{
		ID:        bucket.ID,
		CreatedAt: bucket.CreatedAt.Format(time.DateTime),
		Name:      bucket.Name,
		Capacity:  bucket.Capacity,
	}
}
