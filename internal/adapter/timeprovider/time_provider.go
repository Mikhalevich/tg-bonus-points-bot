package timeprovider

import (
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/manager/orderprocessing"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port"
)

var (
	_ port.TimeProvider            = (*TimeProvider)(nil)
	_ orderprocessing.TimeProvider = (*TimeProvider)(nil)
)

type TimeProvider struct {
}

func New() *TimeProvider {
	return &TimeProvider{}
}

func (t *TimeProvider) Now() time.Time {
	return time.Now()
}
