package port

import (
	"context"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
)

type ButtonRepository interface {
	SetButton(ctx context.Context, b button.Button) (button.InlineKeyboardButton, error)
	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
	SetButtonRows(ctx context.Context, rows ...button.ButtonRow) ([]button.InlineKeyboardButtonRow, error)
	IsNotFoundError(err error) bool
}
