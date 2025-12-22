package port

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
)

type ButtonRepositoryWriter interface {
	SetButton(ctx context.Context, b button.Button) (button.InlineKeyboardButton, error)
	SetButtonRows(ctx context.Context, rows ...button.ButtonRow) ([]button.InlineKeyboardButtonRow, error)
}

type ButtonRepositoryReader interface {
	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
	IsNotFoundError(err error) bool
}
