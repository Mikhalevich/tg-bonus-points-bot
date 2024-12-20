package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

func (o *Order) MakeOrder(ctx context.Context, info msginfo.Info) error {
	orderID := uuid.NewString()

	png, err := o.qrCode.GeneratePNG(orderID)
	if err != nil {
		return fmt.Errorf("qrcode generate png: %w", err)
	}

	if err := o.sender.SendPNG(ctx, info.ChatID, png); err != nil {
		return fmt.Errorf("send png: %w", err)
	}

	return nil
}
