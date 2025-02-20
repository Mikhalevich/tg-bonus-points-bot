package timeprovider

import (
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
)

var _ port.TimeProvider = (*TimeProvider)(nil)

type TimeProvider struct {
}

func New() *TimeProvider {
	return &TimeProvider{}
}

func (t *TimeProvider) Now() time.Time {
	return time.Now()
}
