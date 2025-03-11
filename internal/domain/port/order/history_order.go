package order

import (
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
)

type HistoryOrder struct {
	ID           ID
	SerialNumber int
	Status       Status
	CurrencyID   currency.ID
	CreatedAt    time.Time
	TotalPrice   int
}
