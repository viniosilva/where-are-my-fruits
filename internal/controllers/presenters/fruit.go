package presenters

import (
	"github.com/shopspring/decimal"
)

type CreateFruitReq struct {
	Name      string          `json:"name" example:"Orange"`
	Price     decimal.Decimal `json:"price" example:"1.99"`
	ExpiresIn string          `json:"expires_in" example:"1m"`
	BucketID  *int64          `json:"bucket_id" example:"1"`
}

type FruitRes struct {
	ID        int64  `json:"id" example:"1"`
	CreatedAt string `json:"created_at" example:"2000-12-31 23:59:59"`
	DeletedAt string `json:"deleted_at,omitempty" example:"2000-12-31 23:59:59"`
	BucketID  *int64 `json:"bucket_id,omitempty" example:"1"`

	Name      string          `json:"name" example:"Orange"`
	Price     decimal.Decimal `json:"price" example:"1.99"`
	ExpiresAt string          `json:"expires_at" example:"1m"`
}
