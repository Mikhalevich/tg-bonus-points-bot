package model

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/store"
)

type Store struct {
	ID                int    `db:"id"`
	Description       string `db:"description"`
	DefaultCurrencyID int    `db:"default_currency_id"`
}

func (s *Store) ToPortStore() *store.Store {
	return &store.Store{
		ID:                store.IDFromInt(s.ID),
		Description:       s.Description,
		DefaultCurrencyID: currency.IDFromInt(s.DefaultCurrencyID),
	}
}
