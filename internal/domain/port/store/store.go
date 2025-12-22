package store

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/internal/id"
)

type ID struct {
	id.IntID
}

func IDFromInt(i int) ID {
	return ID{
		IntID: id.IntIDFromInt(i),
	}
}

type Store struct {
	ID                ID
	Description       string
	DefaultCurrencyID currency.ID
	Schedule          Schedule
}
