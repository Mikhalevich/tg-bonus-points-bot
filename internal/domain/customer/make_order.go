package customer

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *Customer) MakeOrder(ctx context.Context, info msginfo.Info) error {
	input := port.CreateOrderInput{
		ChatID:              info.ChatID,
		Status:              order.StatusCreated,
		StatusOperationTime: time.Now(),
		VerificationCode:    generateVerificationCode(),
	}

	id, err := c.repository.CreateOrder(ctx, input)

	if err != nil {
		if c.repository.IsAlreadyExistsError(err) {
			c.sender.ReplyText(ctx, info.ChatID, info.MessageID,
				"You have active order already")
			return nil
		}

		return fmt.Errorf("repository create order: %w", err)
	}

	png, err := c.qrCode.GeneratePNG(id.String())
	if err != nil {
		return fmt.Errorf("qrcode generate png: %w", err)
	}

	orderInfo := formatOrder(&order.Order{
		ID:               id,
		Status:           input.Status,
		VerificationCode: input.VerificationCode,
		Timeline: []order.StatusTime{
			{
				Status: input.Status,
				Time:   input.StatusOperationTime,
			},
		},
	}, c.sender.EscapeMarkdown)

	if err := c.sender.SendPNGMarkdown(ctx, info.ChatID, orderInfo, png); err != nil {
		return fmt.Errorf("send png: %w", err)
	}

	return nil
}

func generateVerificationCode() string {
	//nolint:gosec
	return fmt.Sprintf("%03d", rand.Intn(1000))
}
