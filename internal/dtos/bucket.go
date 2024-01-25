package dtos

type CreateBucketDto struct {
	Name     string `validate:"required,gt=0,lte=128"`
	Capacity int    `validate:"required,gt=0"`
}
