package helpers

import "time"

type TimeImpl struct{}

func (impl *TimeImpl) Now() time.Time {
	return time.Now()
}
