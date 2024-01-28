package presenters

import "github.com/shopspring/decimal"

type CreateBucketReq struct {
	Name     string `json:"name" example:"A"`
	Capacity int    `json:"capacity" example:"10"`
}

type BucketRes struct {
	ID        int64  `json:"id" example:"1"`
	CreatedAt string `json:"created_at" example:"2000-12-31 23:59:59"`
	DeletedAt string `json:"deleted_at,omitempty" example:"2000-12-31 23:59:59"`

	Name     string `json:"name" example:"A"`
	Capacity int    `json:"capacity" example:"10"`
}

type BucketFruitsRes struct {
	ID          int64           `json:"id" example:"1"`
	Name        string          `json:"name" example:"A"`
	Capacity    int             `json:"capacity" example:"10"`
	TotalFruits int64           `json:"total_fruit" example:"5"`
	TotalPrice  decimal.Decimal `json:"total_price" example:"23.54"`
	Percent     string          `json:"percent" example:"50%"`
}

type BucketsFruitsRes struct {
	Data []BucketFruitsRes `json:"data"`
}
