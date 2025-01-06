package customer

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/flag"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) MakeOrder(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID) error {
	input := port.CreateOrderInput{
		ChatID:              chatID,
		Status:              order.StatusAssembling,
		StatusOperationTime: time.Now(),
		VerificationCode:    generateVerificationCode(),
	}

	id, err := c.repository.CreateOrder(ctx, input)

	if err != nil {
		if c.repository.IsAlreadyExistsError(err) {
			c.sender.ReplyText(ctx, chatID, messageID,
				"You have active order already")
			return nil
		}

		return fmt.Errorf("repository create order: %w", err)
	}

	categories, err := c.repository.GetCategoryProducts(ctx, product.Filter{
		Products: flag.Enabled,
		Category: flag.Enabled,
	})
	if err != nil {
		return fmt.Errorf("get products: %w", err)
	}

	buttons := make([]button.InlineKeyboardButton, 0, len(categories)+2)

	for _, v := range categories {
		b, err := c.makeInlineKeyboardButton(ctx, button.CancelOrder(chatID, id), v.Title)
		if err != nil {
			return fmt.Errorf("category order button: %w", err)
		}

		buttons = append(buttons, b)
	}

	cancelBtn, err := c.makeInlineKeyboardButton(ctx, button.CancelOrder(chatID, id), "Cancel")
	if err != nil {
		return fmt.Errorf("cancel order button: %w", err)
	}

	buttons = append(buttons, cancelBtn)

	confirmBtn, err := c.makeInlineKeyboardButton(ctx, button.CancelOrder(chatID, id), "Confirm")
	if err != nil {
		return fmt.Errorf("confirm order button: %w", err)
	}

	buttons = append(buttons, confirmBtn)

	c.sender.ReplyText(ctx, chatID, messageID, "Choose category", buttons...)

	return nil
}

func generateVerificationCode() string {
	//nolint:gosec
	return fmt.Sprintf("%03d", rand.Intn(1000))
}
