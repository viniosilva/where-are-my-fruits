package helpers

import "time"

//go:generate mockgen -source=./time.go -destination=../../mocks/helper_time_mocks.go -package=mocks
type Time interface {
	Now() time.Time
}

type TimeImpl struct{}

func (impl *TimeImpl) Now() time.Time {
	return time.Now()
}
