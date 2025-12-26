package timeprovider

import (
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/cartprocessing"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/orderaction"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/orderpayment"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/manager/orderprocessing"
)

var (
	_ cartprocessing.TimeProvider  = (*TimeProvider)(nil)
	_ orderprocessing.TimeProvider = (*TimeProvider)(nil)
	_ orderaction.TimeProvider     = (*TimeProvider)(nil)
	_ orderpayment.TimeProvider    = (*TimeProvider)(nil)
)

type TimeProvider struct {
}

func New() *TimeProvider {
	return &TimeProvider{}
}

func (t *TimeProvider) Now() time.Time {
	return time.Now()
}
