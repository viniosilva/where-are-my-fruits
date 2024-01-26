package repositories

import (
	"context"
	"database/sql"

	"github.com/viniosilva/where-are-my-fruits/internal/helpers"
	"gorm.io/gorm"
)

var _time helpers.Time

func init() {
	_time = &helpers.TimeImpl{}
}

//go:generate mockgen -source=./interfaces.go -destination=../../mocks/repositories_mocks.go -package=mocks
type DB interface {
	Create(value interface{}) (tx *gorm.DB)
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
}

type SQL interface {
	PingContext(ctx context.Context) error
}
