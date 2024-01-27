package services

import (
	"context"

	"github.com/viniosilva/where-are-my-fruits/internal/helpers"
)

var _time helpers.Time

func init() {
	_time = &helpers.TimeImpl{}
}

//go:generate mockgen -source=./interfaces.go -destination=../../mocks/services_mocks.go -package=mocks
type Logger interface {
	Warn(args ...interface{})
	Error(args ...interface{})
}

type Validate interface {
	Struct(s interface{}) error
}

type DB interface{}

type SQL interface {
	PingContext(ctx context.Context) error
}
