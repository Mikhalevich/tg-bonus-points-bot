package order

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (o *Order) MakeOrder(ctx context.Context, info msginfo.Info) error {
	verifivcationCode := generateVerificationCode()

	id, err := o.repository.CreateOrder(
		ctx,
		port.CreateOrderInput{
			ChatID:              info.ChatID,
			Status:              order.StatusCreated,
			StatusOperationTime: time.Now(),
			VerificationCode:    verifivcationCode,
		})

	if err != nil {
		if o.repository.IsAlreadyExistsError(err) {
			o.sender.ReplyText(ctx, info.ChatID, info.MessageID,
				"You have active order already")
			return nil
		}

		return fmt.Errorf("repository create order: %w", err)
	}

	o.sender.ReplyTextMarkdown(ctx, info.ChatID, info.MessageID,
		fmt.Sprintf("order id: *%s*\n verification code: *%s*", id.String(), verifivcationCode))

	png, err := o.qrCode.GeneratePNG(id.String())
	if err != nil {
		return fmt.Errorf("qrcode generate png: %w", err)
	}

	if err := o.sender.SendPNG(ctx, info.ChatID, png); err != nil {
		return fmt.Errorf("send png: %w", err)
	}

	return nil
}

func generateVerificationCode() string {
	//nolint:gosec
	return fmt.Sprintf("%03d", rand.Intn(1000))
}
