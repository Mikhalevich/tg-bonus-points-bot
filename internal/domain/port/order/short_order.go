package order

import (
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
)

type ShortOrder struct {
	ID               ID
	Status           Status
	VerificationCode string
	CurrencyID       currency.ID
	CreatedAt        time.Time
	TotalPrice       int
}
