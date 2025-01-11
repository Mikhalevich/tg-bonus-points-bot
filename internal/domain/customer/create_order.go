package customer

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *Customer) CreateOrder(ctx context.Context, info msginfo.Info) error {
	input := port.CreateOrderInput{
		ChatID:              info.ChatID,
		Status:              order.StatusConfirmed,
		StatusOperationTime: time.Now(),
		VerificationCode:    generateVerificationCode(),
	}

	id, err := c.repository.CreateOrder(ctx, input)

	if err != nil {
		if c.repository.IsAlreadyExistsError(err) {
			c.sender.ReplyText(ctx, info.ChatID, info.MessageID, message.AlreadyHasActiveOrder())
			return nil
		}

		return fmt.Errorf("repository create order: %w", err)
	}

	if err := c.cart.Clear(ctx, info.ChatID); err != nil {
		return fmt.Errorf("clear cart: %w", err)
	}

	cancelBtn, err := c.makeInlineKeyboardButton(ctx, button.CancelOrder(info.ChatID, id), message.Cancel())
	if err != nil {
		return fmt.Errorf("cancel order button: %w", err)
	}

	png, err := c.qrCode.GeneratePNG(id.String())
	if err != nil {
		return fmt.Errorf("qrcode generate png: %w", err)
	}

	if err := c.sender.SendPNGMarkdown(
		ctx,
		info.ChatID,
		formatOrder(&order.Order{
			ID:               id,
			ChatID:           info.ChatID,
			Status:           input.Status,
			VerificationCode: input.VerificationCode,
		}, c.sender.EscapeMarkdown),
		png,
		button.Row(cancelBtn),
	); err != nil {
		return fmt.Errorf("send png: %w", err)
	}

	c.sender.DeleteMessage(ctx, info.ChatID, info.MessageID)

	return nil
}

func generateVerificationCode() string {
	//nolint:gosec
	return fmt.Sprintf("%03d", rand.Intn(1000))
}
