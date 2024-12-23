package order

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (o *Order) MakeOrder(ctx context.Context, info msginfo.Info) error {
	id, err := o.repository.CreateOrder(ctx, port.CreateOrderInput{
		ChatID:              info.ChatID,
		Status:              order.StatusCreated,
		StatusOperationTime: time.Now(),
		VerificationCode:    "123",
	})

	if err != nil {
		return fmt.Errorf("repository create order: %w", err)
	}

	png, err := o.qrCode.GeneratePNG(id.String())
	if err != nil {
		return fmt.Errorf("qrcode generate png: %w", err)
	}

	if err := o.sender.SendPNG(ctx, info.ChatID, png); err != nil {
		return fmt.Errorf("send png: %w", err)
	}

	return nil
}
