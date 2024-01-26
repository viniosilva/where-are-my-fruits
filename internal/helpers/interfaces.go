package helpers

import "time"

//go:generate mockgen -source=./interfaces.go -destination=../../mocks/helpers_mocks.go -package=mocks
type Time interface {
	Now() time.Time
}
