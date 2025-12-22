package order

import (
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
)

type HistoryOrder struct {
	ID           ID
	SerialNumber int
	Status       Status
	CurrencyID   currency.ID
	CreatedAt    time.Time
	TotalPrice   int
}
